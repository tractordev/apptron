package bridge

import (
	"sync"

	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/qtalk-go/rpc"

	"tractor.dev/hostbridge/bridge/app"
	"tractor.dev/hostbridge/bridge/core"
	"tractor.dev/hostbridge/bridge/menu"
	"tractor.dev/hostbridge/bridge/screen"
	"tractor.dev/hostbridge/bridge/shell"
	"tractor.dev/hostbridge/bridge/window"
)

func NewServer() *rpc.Server {
	mux := rpc.NewRespondMux()

	var listeners sync.Map

	mux.Handle("Listen", rpc.HandlerFunc(func(r rpc.Responder, c *rpc.Call) {
		c.Receive(nil)
		r.Continue(nil)
		listeners.Store(r, struct{}{})
	}))

	core.EventHandler = func(event core.Event) {
		if event.Type > 0 {
			listeners.Range(func(v, _ interface{}) bool {
				r := v.(rpc.Responder)
				if err := r.Send(event); err != nil {
					//log.Println(err)
					listeners.Delete(v)
				}
				return true
			})

		}
	}

	mux.Handle("Shutdown", rpc.HandlerFunc(func(r rpc.Responder, c *rpc.Call) {
		core.Quit()
	}))

	mux.Handle("window", fn.HandlerFrom(window.Module))
	mux.Handle("menu", fn.HandlerFrom(menu.Module))
	mux.Handle("app", fn.HandlerFrom(app.Module))
	mux.Handle("screen", fn.HandlerFrom(screen.Module))
	mux.Handle("shell", fn.HandlerFrom(shell.Module))

	return &rpc.Server{
		Codec:   codec.JSONCodec{},
		Handler: mux,
	}
}
