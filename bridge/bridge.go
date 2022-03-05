package bridge

import (
	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/qtalk-go/rpc"

	"github.com/progrium/hostbridge/bridge/app"
	"github.com/progrium/hostbridge/bridge/core"
	"github.com/progrium/hostbridge/bridge/menu"
	"github.com/progrium/hostbridge/bridge/screen"
	"github.com/progrium/hostbridge/bridge/shell"
	"github.com/progrium/hostbridge/bridge/window"
)

func NewServer() *rpc.Server {
	mux := rpc.NewRespondMux()

	mux.Handle("Listen", rpc.HandlerFunc(func(r rpc.Responder, c *rpc.Call) {
		c.Receive(nil)
		r.Continue(nil)
		core.EventHandler = func(event core.Event) {
			if event.Type > 0 {
				r.Send(event)
			}
		}
	}))

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
