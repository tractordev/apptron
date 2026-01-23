//go:build js && wasm

package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"syscall/js"
	"time"

	"github.com/hugelgupf/p9/p9"
	"github.com/u-root/uio/ulog"
	"tractor.dev/toolkit-go/engine/cli"
	"tractor.dev/wanix"
	"tractor.dev/wanix/fs"
	"tractor.dev/wanix/fs/cowfs"
	"tractor.dev/wanix/fs/fskit"
	"tractor.dev/wanix/fs/httpfs"
	"tractor.dev/wanix/fs/memfs"
	"tractor.dev/wanix/fs/p9kit"
	"tractor.dev/wanix/fs/syncfs"
	"tractor.dev/wanix/fs/tarfs"
	"tractor.dev/wanix/vfs/pipe"
	"tractor.dev/wanix/vfs/ramfs"
	"tractor.dev/wanix/vm"
	"tractor.dev/wanix/vm/v86/virtio9p"
	"tractor.dev/wanix/web"
	"tractor.dev/wanix/web/api"
	"tractor.dev/wanix/web/idbfs"
	"tractor.dev/wanix/web/jsutil"
	wanixruntime "tractor.dev/wanix/web/runtime"
)

var Version string

// todo: centralize or make based on jwt claims
// there is also admins defined in the cloudflare worker
var admins = []string{"progrium"}

func updateLoader(text string) {
	loader := js.Global().Get("document").Call("getElementById", "loader")
	if loader.IsUndefined() {
		return
	}
	loader.Call("querySelector", "p").Set("textContent", text)
}

func main() {
	mainStart := time.Now()
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
		if !user.IsNull() {
			username = user.Get("username").String()
			userID = user.Get("user_id").String()
		}
	}

	envUUID := ""
	envName := ""
	envOwner := ""
	isOwner := false
	env := apptronCfg.Get("env")
	if !env.IsUndefined() {
		envUUID = env.Get("uuid").String()
		envName = env.Get("name").String()
		envOwner = env.Get("owner").String()
		isOwner = envOwner == userID
	}

	if username == "" {
		log.Printf("starting apptron wanix for anonymous, env %s\n", envName)
	} else {
		log.Printf("starting apptron wanix for user %s, env %s\n", username, envName)
	}

	inst := wanixruntime.Instance()

	k := wanix.New()
	k.AddModule("#web", web.New(k))
	k.AddModule("#vm", vm.New())
	k.AddModule("#pipe", &pipe.Allocator{})
	k.AddModule("#commands", &pipe.Allocator{})
	k.AddModule("#|", &pipe.Allocator{}) // alias for #pipe
	k.AddModule("#ramfs", &ramfs.Allocator{})

	root, err := k.NewRoot()
	if err != nil {
		log.Fatal(err)
	}

	// setup 9p early to ensure hook is ready for vm
	debug9p := inst.Get("config").Get("debug9p")
	if debug9p.IsUndefined() {
		debug9p = js.ValueOf(false)
	}
	// cached9p := metacache.New(root.Namespace())
	// cached9p.SetLogger(slogger.New(slog.LevelDebug))
	// run9p := virtio9p.Setup(cached9p, inst, debug9p.Bool())
	run9p := virtio9p.Setup(root.Namespace(), inst, debug9p.Bool())

	// experimental 9p server over ... 9p.
	// this ends up being much slower than using a messagechannel for the 9p connection.
	p9fs := fskit.OpenFunc(func(ctx context.Context, path string) (fs.File, error) {
		p1, p2 := pipe.New(false)
		var o []p9.ServerOpt
		if true {
			o = append(o, p9.WithServerLogger(ulog.Log))
		}
		srv := p9.NewServer(p9kit.Attacher(root.Namespace(), p9kit.WithMemAttrStore()), o...)
		go srv.Handle(p2, p2)
		return fskit.NewStreamFile(p1, p1, p1, "<9p-connection>"), nil
	})
	if err := root.Namespace().Bind(p9fs, ".", "#9psrv"); err != nil {
		log.Fatal(err)
	}

	// setup root bindings
	rootBindings := []struct {
		dst string
		src string
	}{
		{"#task", "task"},
		{"#cap", "cap"},
		{"#web", "web"},
		{"#vm", "vm"},
		{"#|", "#console"}, // todo: is this not allocating?
	}
	for _, b := range rootBindings {
		if err := root.Bind(b.dst, b.src); err != nil {
			log.Fatal(err)
		}
	}

	// set up primary port responder and factory function
	go api.PortResponder(inst.Call("_portConn", inst.Get("_sys").Get("port1")), root)
	inst.Set("createPort", js.FuncOf(func(this js.Value, args []js.Value) any {
		ch := js.Global().Get("MessageChannel").New()
		go api.PortResponder(inst.Call("_portConn", ch.Get("port1")), root)
		return ch.Get("port2")
	}))

	// load bundle
	startTime := time.Now()
	bundleBytes := inst.Get("_bundle")
	if bundleBytes.IsUndefined() {
		log.Fatal("bundle not found")
	}
	if bundleBytes.InstanceOf(js.Global().Get("Promise")) {
		inst.Set("_bundle", jsutil.Await(bundleBytes))
		bundleBytes = inst.Get("_bundle")
	}
	jsBuf := js.Global().Get("Uint8Array").New(bundleBytes)
	b := make([]byte, jsBuf.Length())
	js.CopyBytesToGo(b, jsBuf)
	inst.Set("_bundle", js.Undefined())
	buf := bytes.NewBuffer(b)
	bundleFS := tarfs.From(tar.NewReader(buf))
	if err := root.Namespace().Bind(bundleFS, ".", "#bundle"); err != nil {
		log.Fatal(err)
	}
	if err := root.Bind("#bundle", "bundle"); err != nil {
		log.Fatal(err)
	}
	bundleTime := time.Since(startTime)
	log.Printf("bundle loaded in %v\n", bundleTime)

	// IDBFS is still origin-private if not exactly OPFS.
	// Not only does it work in older Safari, but it's 50% faster than OPFS.
	opfs := idbfs.New("apptron-rev1")
	// opfs.SetLogger(log.Default())
	if err := root.Namespace().Bind(opfs, ".", "web/idbfs/apptron"); err != nil {
		log.Fatal(err)
	}

	// load environment base and scratch
	startTime = time.Now()
	envBase, err := fs.Sub(bundleFS, "rootfs")
	if err != nil {
		log.Fatal(err)
	}
	var envScratch fs.FS = memfs.New()
	root.Namespace().Bind(envScratch, ".", "#scratch")
	if !env.IsUndefined() {
		if rootExists, _ := fs.DirExists(opfs, fmt.Sprintf("env/%s/root", envUUID)); rootExists {
			updateLoader("Loading custom environment...")
			log.Println("using custom env base")
			envBase, err = fs.Sub(opfs, fmt.Sprintf("env/%s/root", envUUID))
			if err != nil {
				log.Fatal(err)
			}
		}
		if overlayExists, _ := fs.DirExists(opfs, fmt.Sprintf("env/%s/overlay", envUUID)); overlayExists {
			updateLoader("Loading custom environment...")
			log.Println("using env overlay")
			envOverlay, err := fs.Sub(opfs, fmt.Sprintf("env/%s/overlay", envUUID))
			if err != nil {
				log.Fatal(err)
			}
			envBase = fskit.UnionFS{
				envBase,
				envOverlay,
			}
		}
	}
	root.Namespace().Bind(&cowfs.FS{Base: envBase, Overlay: envScratch}, ".", "#env")
	root.Namespace().Bind(fskit.RawNode([]byte(Version+"\n")), ".", "#version")
	envTime := time.Since(startTime)
	log.Printf("env (cow) loaded in %v\n", envTime)

	// old way of loading env:
	// startTime = time.Now()
	// envRoot := memfs.New()
	// var envRootExists bool
	// if !env.IsUndefined() {
	// 	envRootExists, _ = fs.DirExists(opfs, fmt.Sprintf("env/%s/root", envUUID))
	// }
	// if envRootExists {
	// 	updateLoader("Loading custom environment...")
	// 	log.Println("using custom env root")
	// 	if err := fs.CopyFS(opfs, fmt.Sprintf("env/%s/root", envUUID), envRoot, "."); err != nil {
	// 		log.Fatal(err)
	// 	}
	// } else {
	// 	if err := fs.CopyFS(bundleFS, "rootfs", envRoot, "."); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// root.Namespace().Bind(envRoot, ".", "#env")
	// envTime = time.Since(startTime)
	// log.Printf("env loaded in %v\n", envTime)

	// setup vm
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
		{"#ramfs", fmt.Sprintf("vm/%s/fsys/#ramfs", vm)},
		{"#pipe", fmt.Sprintf("vm/%s/fsys/#pipe", vm)},
		{"#|", fmt.Sprintf("vm/%s/fsys/#|", vm)},
		{".", fmt.Sprintf("vm/%s/fsys", vm)},
		{"#env", fmt.Sprintf("vm/%s/fsys", vm)},
	}
	for _, b := range vmBindings {
		if err := root.Bind(b.dst, b.src); err != nil {
			log.Fatal(err)
		}
	}

	profile := []string{
		fmt.Sprintf("export USER=%s", username),
	}
	if username != "" {
		profile = append(profile, fmt.Sprintf("export HOME=/home/%s", username))
	}
	if !env.IsUndefined() {
		profile = append(profile, fmt.Sprintf("export ENV_NAME=%s", envName))
		profile = append(profile, fmt.Sprintf("export ENV_UUID=%s", envUUID))
		if username == "" {
			profile = append(profile, fmt.Sprintf("export HOME=/home/%s", envName))
		}
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
		fmt.Sprintf("rootflags=trans=virtio,version=9p2000.L,aname=vm/%s/fsys,cache=loose,msize=131072", vm),
		"mem=1008M",
		"memmap=16M$1008M",
	}
	ctlcmd := []string{
		"start",
		"-m", "1G",
		"-append", fmt.Sprintf("'%s'", strings.Join(cmdline, " ")),
	}
	if !inst.Get("config").Get("network").IsUndefined() {
		ctlcmd = append(ctlcmd, "-netdev")
		ctlcmd = append(ctlcmd, fmt.Sprintf("user,type=virtio,relay_url=%s", inst.Get("config").Get("network").String()))
	}

	// boot vm as early as possible
	log.Println("booting vm")
	if err := fs.WriteFile(root.Namespace(), fmt.Sprintf("vm/%s/ctl", vm), []byte(strings.Join(ctlcmd, " ")), 0755); err != nil {
		log.Fatal(err)
	}

	// go func() {
	// 	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// 	defer cancel()
	// 	if err := fsutil.WaitFor(ctx, root.Namespace(), fmt.Sprintf("vm/%s/fsys/run/shm9p.lock", vm), true); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	shmpipe, err := fs.OpenFile(root.Namespace(), fmt.Sprintf("vm/%s/shmpipe0", vm), os.O_RDWR, 0)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer shmpipe.Close()

	// 	vmfs, err := p9kit.ClientFS(&rwcConn{rwc: fs.DefaultFile{File: shmpipe}}, "/", p9.WithMessageSize(512*1024))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	// sanity check
	// 	entries, err := fs.ReadDir(vmfs, ".")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	if len(entries) == 0 {
	// 		log.Fatal("vmfs is empty, this should not happen")
	// 	}
	// 	if err := root.Namespace().Bind(vmfs, ".", fmt.Sprintf("vm/%s/9proot", vm)); err != nil {
	// 		log.Fatal(err)
	// 	}

	// }()

	// setup control file
	setupBundle := func(name string, rw bool) {
		startTime := time.Now()
		bundle := jsutil.Await(inst.Call("_getBundle", name))
		if bundle.IsUndefined() {
			log.Printf("bundle %s not found\n", name)
			return
		}
		jsBuf := js.Global().Get("Uint8Array").New(bundle)
		b := make([]byte, jsBuf.Length())
		js.CopyBytesToGo(b, jsBuf)
		buf := bytes.NewBuffer(b)
		var fsys fs.FS
		if rw {
			runtime.GC()
			rwfs := memfs.New()
			if err := fs.CopyFS(tarfs.From(tar.NewReader(buf)), ".", rwfs, "."); err != nil {
				log.Fatal(err)
			}
			buf.Reset()
			fsys = rwfs
		} else {
			fsys = tarfs.From(tar.NewReader(buf))
		}
		mountname := filepath.Base(name)
		// Remove any suffix after a dot (including the dot)
		if dot := strings.IndexByte(mountname, '.'); dot != -1 {
			mountname = mountname[:dot]
		}
		if err := root.Namespace().Bind(fsys, ".", "#"+mountname); err != nil {
			log.Fatal(err)
		}
		bundleTime := time.Since(startTime)
		runtime.GC() // free memory
		log.Printf("%s bundle loaded in %v\n", mountname, bundleTime)
	}

	if err := root.Namespace().Bind(wanix.ControlFile(&cli.Command{
		Usage: "ctl",
		Run: func(_ *cli.Context, args []string) {
			log.Println("ctl:", args)
			switch args[0] {
			case "cmd":
				if len(args) < 2 {
					fmt.Println("usage: cmd <cmd>")
					return
				}
				if err := fs.AppendFile(root.Namespace(), "#commands/data", []byte(strings.Join(args[1:], " "))); err != nil {
					log.Fatal(err)
				}

			case "bind":
				if len(args) < 2 {
					fmt.Println("usage: bind <oldname> <newname>")
					return
				}
				if err := root.Bind(args[1], args[2]); err != nil {
					log.Fatal(err)
				}

			case "reload":
				js.Global().Get("location").Call("reload")

			case "bundle":
				if len(args) < 1 {
					fmt.Println("usage: bundle <name>")
					return
				}
				rw := false
				if len(args) > 2 {
					rw = args[2] == "rw"
				}
				setupBundle(fmt.Sprintf("bundles/%s.tar.br", args[1]), rw)
			case "cp":
				if len(args) < 2 {
					fmt.Println("usage: cp <src> <dst>")
					return
				}
				if err := fs.CopyAll(root.Namespace(), args[1], args[2]); err != nil {
					log.Fatal(err)
				}
			case "sync":
				if len(args) < 2 {
					fmt.Println("usage: sync <src> <dst>")
					return
				}
				src := args[1]
				dst := args[2]
				if err := fs.CopyAll(root.Namespace(), src, dst); err != nil {
					log.Fatal(err)
				}
				var toRemove []string
				isRemovedSubpath := func(candidate string) bool {
					for _, parent := range toRemove {
						if parent == candidate {
							return true
						}
						if strings.HasPrefix(candidate, parent+"/") {
							return true
						}
					}
					return false
				}
				fs.WalkDir(root.Namespace(), dst, func(path string, info fs.DirEntry, err error) error {
					if err != nil {
						return nil // skip errors
					}
					relPath, err := filepath.Rel(dst, path)
					if err != nil || relPath == "." {
						return nil // skip root itself or error
					}
					srcPath := filepath.Join(src, relPath)
					ok, err := fs.Exists(root.Namespace(), srcPath)
					if err != nil || !ok {
						fullPath := filepath.Join(dst, relPath)
						if !isRemovedSubpath(fullPath) {
							toRemove = append(toRemove, fullPath)
						}
					}
					return nil
				})
				for _, p := range toRemove {
					if err := fs.RemoveAll(root.Namespace(), p); err != nil {
						log.Fatal(err)
					}
				}
			}
		},
	}), ".", "ctl"); err != nil {
		log.Fatal(err)
	}

	// load environment buildfs
	var buildScratch fs.FS = memfs.New()
	root.Namespace().Bind(buildScratch, ".", "#envbuild")
	buildBase, err := fs.Sub(bundleFS, "rootfs")
	if err != nil {
		log.Fatal(err)
	}
	if err := root.Namespace().Bind(&cowfs.FS{Base: buildBase, Overlay: buildScratch}, ".", fmt.Sprintf("vm/%s/fsys/apptron/.buildroot", vm)); err != nil {
		log.Fatal(err)
	}

	updateLoader("Syncing filesystem...")

	// setup user fs
	if username != "" {
		log.Println("setting up user fs")
		startTime = time.Now()
		remoteHomeFS := httpfs.New(fmt.Sprintf("%s/data/usr/%s", origin.String(), userID), nil)
		// if err := fs.MkdirAll(remoteHomeFS, ".", 0755); err != nil {
		// 	log.Fatal(err)
		// }
		if err := fs.MkdirAll(opfs, path.Join("usr", userID), 0755); err != nil {
			log.Fatal(err)
		}
		localHomeFS, err := fs.Sub(opfs, path.Join("usr", userID))
		if err != nil {
			log.Fatal(err)
		}
		sfs := syncfs.New(localHomeFS, remoteHomeFS, 5*time.Second)

		if err := sfs.Sync(); err != nil {
			log.Printf("err syncing: %v\n", err)
		}

		if err := root.Namespace().Bind(sfs, ".", fmt.Sprintf("home/%s", username)); err != nil {
			log.Fatal(err)
		}
		log.Printf("user fs ready in %v\n", time.Since(startTime))
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		// setup datafs for admins
		isAdmin := slices.Contains(admins, username)
		if isAdmin {
			datafs := httpfs.New(fmt.Sprintf("%s/data", origin.String()), nil)
			// datafs.SetLogger(slogger.New(slog.LevelDebug))
			datafs.Ignore("MAILPATH") // annoying busybox thing
			cachedDatafs := httpfs.NewCacher(datafs)
			go func() {
				if _, _, err := cachedDatafs.PullNode(context.Background(), ".", true); err != nil {
					log.Printf("err pulling datafs: %v\n", err)
				}
			}()
			if err := root.Namespace().Bind(cachedDatafs, ".", "root/data"); err != nil {
				log.Fatal(err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// setup project fs
		if !env.IsUndefined() {
			startTime = time.Now()

			var localProjectFS fs.FS
			var err error
			if isOwner {
				log.Println("setting up project fs")
				if err := fs.MkdirAll(opfs, path.Join("env", envUUID, "project"), 0755); err != nil {
					log.Fatal(err)
				}
				localProjectFS, err = fs.Sub(opfs, path.Join("env", envUUID, "project"))
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Println("setting up temp project fs (not owner)")
				localProjectFS = memfs.New()
			}
			remoteProjectFS := httpfs.New(fmt.Sprintf("%s/data/env/%s/project", origin.String(), envUUID), nil)

			sfs := syncfs.New(localProjectFS, remoteProjectFS, 5*time.Second)
			if err := sfs.Sync(); err != nil {
				log.Fatalf("err syncing: %v %v\n", err, envUUID)
			}

			if err := root.Namespace().Bind(sfs, ".", "project"); err != nil {
				log.Fatal(err)
			}
			if err := root.Bind("project", filepath.Join("home", username, envName)); err != nil {
				log.Fatal(err)
			}
			log.Printf("project fs ready in %v\n", time.Since(startTime))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// setup public fs
		if isOwner {
			log.Println("setting up public fs")
			remotePublicFS := httpfs.New(fmt.Sprintf("%s/data/env/%s/public", origin.String(), envUUID), nil)
			if err := fs.MkdirAll(remotePublicFS, ".", 0755); err != nil {
				log.Fatal(err)
			}
			startTime = time.Now()
			if err := fs.MkdirAll(opfs, path.Join("env", envUUID, "public"), 0755); err != nil {
				log.Fatal(err)
			}
			localPublicFS, err := fs.Sub(opfs, path.Join("env", envUUID, "public"))
			if err != nil {
				log.Fatal(err)
			}
			publicSyncFS := syncfs.New(localPublicFS, remotePublicFS, 3*time.Second)
			if err := publicSyncFS.Sync(); err != nil {
				log.Printf("err syncing: %v\n", err)
			}

			if err := root.Namespace().Bind(publicSyncFS, ".", "public"); err != nil {
				log.Fatal(err)
			}
			log.Printf("public fs ready in %v\n", time.Since(startTime))
		}
	}()

	wg.Wait()

	// done
	inst.Call("_wasmReady")
	log.Printf("apptron ready in %v\n", time.Since(mainStart))

	// block on serving 9p
	run9p()
}

// Conn is an adapter that implements net.Conn using an underlying io.ReadWriteCloser.
// LocalAddr/RemoteAddr will be dummy addrs, SetDeadline/Set[Read|Write]Deadline are no-ops.
type rwcConn struct {
	rwc io.ReadWriteCloser
}

func (c *rwcConn) Read(b []byte) (int, error) {
	return c.rwc.Read(b)
}
func (c *rwcConn) Write(b []byte) (int, error) {
	return c.rwc.Write(b)
}
func (c *rwcConn) Close() error {
	return c.rwc.Close()
}
func (c *rwcConn) LocalAddr() (addr net.Addr) {
	return dummyAddr("rwc-local")
}
func (c *rwcConn) RemoteAddr() (addr net.Addr) {
	return dummyAddr("rwc-remote")
}
func (c *rwcConn) SetDeadline(t time.Time) error {
	return nil // not supported
}
func (c *rwcConn) SetReadDeadline(t time.Time) error {
	return nil // not supported
}
func (c *rwcConn) SetWriteDeadline(t time.Time) error {
	return nil // not supported
}

type dummyAddr string

func (a dummyAddr) Network() string { return string(a) }
func (a dummyAddr) String() string  { return string(a) }
