package bridge

import (
	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/qtalk-go/rpc"

	"tractor.dev/hostbridge/bridge/api/app"
	"tractor.dev/hostbridge/bridge/api/menu"
	"tractor.dev/hostbridge/bridge/api/screen"
	"tractor.dev/hostbridge/bridge/api/shell"
	"tractor.dev/hostbridge/bridge/api/window"
	"tractor.dev/hostbridge/bridge/event"
	"tractor.dev/hostbridge/bridge/platform"
)

func NewServer() *rpc.Server {
	mux := rpc.NewRespondMux()

	mux.Handle("Listen", rpc.HandlerFunc(func(r rpc.Responder, c *rpc.Call) {
		c.Receive(nil)
		r.Continue(nil)
		event.Listen(r, func(e event.Event) error {
			return r.Send(e)
		})
	}))

	mux.Handle("Shutdown", rpc.HandlerFunc(func(r rpc.Responder, c *rpc.Call) {
		platform.Terminate()
	}))

	mux.Handle("window", fn.HandlerFrom(window.Module))
	mux.Handle("menu", fn.HandlerFrom(menu.Module))
	mux.Handle("app", fn.HandlerFrom(app.Module))
	mux.Handle("screen", fn.HandlerFrom(screen.Module))
	mux.Handle("shell", fn.HandlerFrom(shell.Module))

	return &rpc.Server{
		Codec: codec.JSONCodec{},
		Handler: rpc.HandlerFunc(func(r rpc.Responder, c *rpc.Call) {
			platform.Dispatch(func() {
				mux.RespondRPC(r, c)
			})
		}),
	}
}
