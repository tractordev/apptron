package apputil

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/qtalk-go/rpc"
	"tractor.dev/hostbridge/client"
)

// Run takes a filesystem and optional userMethods value to start a simple hostbridge
// program that serves the filesystem over HTTP to be used by hostbridge windows. It
// launches a main window pointing at index.html from the filesystem, using meta tags
// via OptionsFromHTML to define options for this main window. It also serves a backend
// using BackendServer that lets you include /-/hostbridge.js which creates a "$host" global
// that exposes the hostbridge API. This API can be extended by the methods on the
// userMethods value, made accessible as callables via "$host.rpc.user".
//
// If the hostbridge binary is included in the filesystem as "./hostbridge", then it will
// use TempCommand to write it to disk and the BRIDGECMD environment variable to use it.
//
// This can be used either as-is or as a reference for your own hostbridge programs, but is
// mainly used by programs produced by the built-in build subcommand.
func Run(fsys fs.FS, userMethods interface{}) {
	if path, cleanup := TempCommand(fsys, "hostbridge"); path != "" {
		os.Setenv("BRIDGECMD", path)
		defer cleanup()
	}

	bridge, err := client.Spawn()
	if err != nil {
		log.Fatal(err)
	}
	defer bridge.Close()

	bridge.OnEvent = func(event client.Event) {
		if event.Type == client.EventClose && event.WindowID == 1 {
			bridge.Close()
		}
	}

	l, err := net.Listen("tcp4", ":0")
	if err != nil {
		log.Fatal(err)
	}
	srv := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/-/") {
				BackendServer(bridge, func(mux *rpc.RespondMux) {
					if userMethods != nil {
						mux.Handle("user", fn.HandlerFrom(userMethods))
					}
				}).ServeHTTP(w, r)
				return
			}
			http.FileServer(http.FS(fsys)).ServeHTTP(w, r)
		}),
	}
	go srv.Serve(l)

	ctx := context.Background()
	_, err = bridge.Window.New(ctx, OptionsFromHTML(fsys, "index.html", "window", client.WindowOptions{
		Size: client.Size{
			Width:  640,
			Height: 480,
		},
		Visible: true,
		URL:     fmt.Sprintf("http://%s/", l.Addr().String()),
	}))
	if err != nil {
		log.Fatal(err)
	}

	bridge.Wait()
}
