//go:build js && wasm

package main

import (
	"archive/tar"
	"bytes"
	"log"
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
	"tractor.dev/wanix/web/runtime"
	"tractor.dev/wanix/web/virtio9p"
)

func main() {
	log.SetFlags(log.Lshortfile)

	log.Println("starting apptron wanix")

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

	// todo: let config define, otherwise default to these
	root.Bind("#task", "task")
	root.Bind("#cap", "cap")
	root.Bind("#web", "web")
	root.Bind("#vm", "vm")
	root.Bind("#|", "#console")

	bundleBytes := inst.Get("_bundle")
	if !bundleBytes.IsUndefined() {
		jsBuf := js.Global().Get("Uint8Array").New(bundleBytes)
		b := make([]byte, jsBuf.Length())
		js.CopyBytesToGo(b, jsBuf)
		buf := bytes.NewBuffer(b)
		bundleFS := tarfs.From(tar.NewReader(buf))

		// ideally we could bind a memfs over bundleFS, but
		// that still doesn't seem to be working yet
		rw := memfs.New()
		if err := fs.CopyFS(bundleFS, ".", rw, "."); err != nil {
			log.Fatal(err)
		}
		root.Namespace().Bind(rw, ".", "#bundle")
		// root.Bind("#bundle", "bundle")
	}

	r2fs := httpfs.New("https://r2fs.proteco.workers.dev/", nil)
	opfs, err := fsa.OPFS("r2fs")
	if err != nil {
		log.Fatal(err)
	}
	sfs := syncfs.New(opfs, r2fs, 3*time.Second)
	go func() {
		if err := sfs.Sync(); err != nil {
			log.Printf("err syncing: %v\n", err)
		}
	}()
	if err := root.Namespace().Bind(sfs, ".", "#data"); err != nil {
		log.Fatal(err)
	}

	inst.Set("createPort", js.FuncOf(func(this js.Value, args []js.Value) any {
		ch := js.Global().Get("MessageChannel").New()
		go api.PortResponder(inst.Call("_portConn", ch.Get("port1")), root)
		return ch.Get("port2")
	}))

	go api.PortResponder(inst.Call("_portConn", inst.Get("_sys").Get("port1")), root)

	inst.Call("_wasmReady")

	virtio9p.Serve(root.Namespace(), inst, false)

}
