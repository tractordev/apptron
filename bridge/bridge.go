package bridge

import (
	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/qtalk-go/rpc"
	"github.com/progrium/qtalk-go/x/cbor/codec"

	"tractor.dev/apptron/bridge/api/app"
	"tractor.dev/apptron/bridge/api/menu"
	"tractor.dev/apptron/bridge/api/shell"
	"tractor.dev/apptron/bridge/api/system"
	"tractor.dev/apptron/bridge/api/window"
	"tractor.dev/apptron/bridge/event"
	"tractor.dev/apptron/bridge/platform"
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
		platform.Terminate(false)
	}))

	mux.Handle("window", fn.HandlerFrom(window.Module))
	mux.Handle("menu", fn.HandlerFrom(menu.Module))
	mux.Handle("app", fn.HandlerFrom(app.Module))
	mux.Handle("system", fn.HandlerFrom(system.Module))
	mux.Handle("shell", fn.HandlerFrom(shell.Module))

	return &rpc.Server{
		Codec: codec.CBORCodec{},
		Handler: rpc.HandlerFunc(func(r rpc.Responder, c *rpc.Call) {
			platform.Dispatch(func() {
				mux.RespondRPC(r, c)
			})
		}),
	}
}
