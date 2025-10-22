//go:build js && wasm

package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"slices"
	"strings"
	"syscall/js"
	"time"

	"tractor.dev/wanix"
	"tractor.dev/wanix/fs"
	"tractor.dev/wanix/fs/httpfs"
	"tractor.dev/wanix/fs/memfs"
	"tractor.dev/wanix/fs/syncfs"
	"tractor.dev/wanix/fs/tarfs"
	"tractor.dev/wanix/vfs/pipe"
	"tractor.dev/wanix/vfs/ramfs"
	"tractor.dev/wanix/vm"
	"tractor.dev/wanix/web"
	"tractor.dev/wanix/web/api"
	"tractor.dev/wanix/web/fsa"
	"tractor.dev/wanix/web/jsutil"
	"tractor.dev/wanix/web/runtime"
	"tractor.dev/wanix/web/virtio9p"
)

// todo: centralize or make based on jwt claims
// there is also admins defined in the cloudflare worker
var admins = []string{"progrium"}

func main() {
	log.SetFlags(log.Lshortfile)

	apptronCfg := js.Global().Get("window").Get("apptron")
	if apptronCfg.IsUndefined() {
		log.Fatal("apptron config not found")
	}
	origin := apptronCfg.Get("origin")
	if origin.IsUndefined() {
		log.Fatal("apptron origin not found")
	}
	username := ""
	userID := ""
	user := apptronCfg.Get("user")
	if !user.IsUndefined() {
		if user.InstanceOf(js.Global().Get("Promise")) {
			apptronCfg.Set("user", jsutil.Await(user))
			user = apptronCfg.Get("user")
		}
		username = user.Get("username").String()
		userID = user.Get("user_id").String()
	}

	envUUID := ""
	envName := ""
	env := apptronCfg.Get("env")
	if !env.IsUndefined() {
		envUUID = env.Get("uuid").String()
		envName = env.Get("name").String()
	}

	log.Printf("starting apptron wanix for user %s, env %s\n", username, envName)

	inst := runtime.Instance()

	k := wanix.New()
	k.AddModule("#web", web.New(k))
	k.AddModule("#vm", vm.New())
	k.AddModule("#pipe", &pipe.Allocator{})
	k.AddModule("#|", &pipe.Allocator{}) // alias for #pipe
	k.AddModule("#ramfs", &ramfs.Allocator{})

	root, err := k.NewRoot()
	if err != nil {
		log.Fatal(err)
	}

	bundleBytes := inst.Get("_bundle")
	if bundleBytes.IsUndefined() {
		log.Fatal("bundle not found")
	}
	jsBuf := js.Global().Get("Uint8Array").New(bundleBytes)
	b := make([]byte, jsBuf.Length())
	js.CopyBytesToGo(b, jsBuf)
	buf := bytes.NewBuffer(b)
	bundleFS := tarfs.From(tar.NewReader(buf))
	if err := root.Namespace().Bind(bundleFS, ".", "#bundle"); err != nil {
		log.Fatal(err)
	}

	// todo: let config define, otherwise default to these
	rootBindings := []struct {
		dst string
		src string
	}{
		{"#task", "task"},
		{"#cap", "cap"},
		{"#web", "web"},
		{"#vm", "vm"},
		{"#|", "#console"},
		{"#bundle", "bundle"},
	}
	for _, b := range rootBindings {
		if err := root.Bind(b.dst, b.src); err != nil {
			log.Fatal(err)
		}
	}

	opfs, err := fsa.OPFS("apptron")
	if err != nil {
		log.Fatal(err)
	}

	envRoot := memfs.New()
	var envRootExists bool
	if !env.IsUndefined() {
		envRootExists, _ = fs.DirExists(opfs, fmt.Sprintf("env/%s/root", envUUID))
	}
	if envRootExists {
		log.Println("using custom env root")
		if err := fs.CopyFS(opfs, fmt.Sprintf("env/%s/root", envUUID), envRoot, "."); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := fs.CopyFS(bundleFS, "rootfs", envRoot, "."); err != nil {
			log.Fatal(err)
		}
	}
	root.Namespace().Bind(envRoot, ".", "#env")

	envBuild := memfs.New()
	if err := fs.CopyFS(bundleFS, "rootfs", envBuild, "."); err != nil {
		log.Fatal(err)
	}
	root.Namespace().Bind(envBuild, ".", "envbuild")

	// opfs, err := fsa.OPFS("apptron")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fs.MkdirAll(opfs, "sys/overlay", 0755)
	// base, err := fs.Sub(opfs, "sys/root")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// overlay, err := fs.Sub(opfs, "sys/overlay")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// cfs := &cowfs.FS{
	// 	Base:    base,
	// 	Overlay: overlay,
	// }
	// if err := cfs.Whiteout(".whiteout"); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := root.Namespace().Bind(cfs, ".", "cow"); err != nil {
	// 	log.Fatal(err)
	// }

	isAdmin := slices.Contains(admins, username)
	if isAdmin {
		datafs := httpfs.New(fmt.Sprintf("%s/data", origin.String()), nil)
		// datafs.SetLogger(slogger.New(slog.LevelDebug))
		datafs.Ignore("MAILPATH")
		cachedDatafs := httpfs.NewCacher(datafs)
		_, _, err := cachedDatafs.PullNode(context.Background(), ".", true)
		if err != nil {
			log.Fatal(err)
		}
		if err := root.Namespace().Bind(cachedDatafs, ".", "root/data"); err != nil {
			log.Fatal(err)
		}
	}

	remoteHomeFS := httpfs.New(fmt.Sprintf("%s/data/usr/%s", origin.String(), userID), nil)
	localHomeFS, err := fs.Sub(opfs, fmt.Sprintf("usr/%s", userID))
	if err != nil {
		log.Fatal(err)
	}
	sfs := syncfs.New(localHomeFS, remoteHomeFS, 3*time.Second)
	log.Println("syncing user fs")
	if err := sfs.Sync(); err != nil {
		log.Printf("err syncing: %v\n", err)
	}
	log.Println("user fs synced")
	if err := root.Namespace().Bind(sfs, ".", fmt.Sprintf("home/%s", username)); err != nil {
		log.Fatal(err)
	}

	if !env.IsUndefined() {
		remoteProjectFS := httpfs.New(fmt.Sprintf("%s/data/env/%s/project", origin.String(), envUUID), nil)
		localProjectFS, err := fs.Sub(opfs, fmt.Sprintf("env/%s/project", envUUID))
		if err != nil {
			log.Fatal(err)
		}
		sfs := syncfs.New(localProjectFS, remoteProjectFS, 3*time.Second)
		log.Println("syncing project fs")
		if err := sfs.Sync(); err != nil {
			log.Printf("err syncing: %v\n", err)
		}
		log.Println("project fs synced")
		if err := root.Namespace().Bind(sfs, ".", "project"); err != nil {
			log.Fatal(err)
		}
		if err := root.Bind("project", filepath.Join("home", username, envName)); err != nil {
			log.Fatal(err)
		}
	}

	inst.Set("createPort", js.FuncOf(func(this js.Value, args []js.Value) any {
		ch := js.Global().Get("MessageChannel").New()
		go api.PortResponder(inst.Call("_portConn", ch.Get("port1")), root)
		return ch.Get("port2")
	}))

	go api.PortResponder(inst.Call("_portConn", inst.Get("_sys").Get("port1")), root)

	// boot bundle

	vmraw, err := fs.ReadFile(root.Namespace(), "vm/new/default")
	if err != nil {
		log.Fatal(err)
	}
	vm := strings.TrimSpace(string(vmraw))

	vmBindings := []struct {
		dst string
		src string
	}{
		{"#console/data1", fmt.Sprintf("vm/%s/ttyS0", vm)},
		{".", fmt.Sprintf("vm/%s/fsys", vm)},
		{"#ramfs", fmt.Sprintf("vm/%s/fsys/#ramfs", vm)},
		{"#pipe", fmt.Sprintf("vm/%s/fsys/#pipe", vm)},
		{"#|", fmt.Sprintf("vm/%s/fsys/#|", vm)},
		{"#env", fmt.Sprintf("vm/%s/fsys", vm)},
	}
	for _, b := range vmBindings {
		if err := root.Bind(b.dst, b.src); err != nil {
			log.Fatal(err)
		}
	}

	profile := []string{
		fmt.Sprintf("export USER=%s", username),
		fmt.Sprintf("export HOME=/home/%s", username),
	}
	if !env.IsUndefined() {
		profile = append(profile, fmt.Sprintf("export ENV_NAME=%s", envName))
		profile = append(profile, fmt.Sprintf("export ENV_UUID=%s", envUUID))
	}
	profile = append(profile, "")
	if err := fs.WriteFile(root.Namespace(), "#env/etc/profile.d/apptron.sh", []byte(strings.Join(profile, "\n")), 0644); err != nil {
		log.Fatal(err)
	}
	cmdline := []string{
		"init=/bin/init",
		"rw",
		"root=host9p",
		"rootfstype=9p",
		fmt.Sprintf("rootflags=trans=virtio,version=9p2000.L,aname=vm/%s/fsys,cache=none", vm),
	}
	ctlcmd := []string{
		"start",
		"-append", fmt.Sprintf("'%s'", strings.Join(cmdline, " ")),
	}
	if !inst.Get("config").Get("network").IsUndefined() {
		ctlcmd = append(ctlcmd, "-netdev")
		ctlcmd = append(ctlcmd, fmt.Sprintf("user,type=virtio,relay_url=%s", inst.Get("config").Get("network").String()))
	}
	if err := fs.WriteFile(root.Namespace(), fmt.Sprintf("vm/%s/ctl", vm), []byte(strings.Join(ctlcmd, " ")), 0755); err != nil {
		log.Fatal(err)
	}

	inst.Call("_wasmReady")
	log.Println("apptron ready")

	virtio9p.Serve(root.Namespace(), inst, false)
}
