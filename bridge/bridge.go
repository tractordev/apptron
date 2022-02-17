package bridge

import (
	"unsafe"

	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/qtalk-go/rpc"

	"github.com/progrium/hostbridge/bridge/app"
	"github.com/progrium/hostbridge/bridge/window"
)

type ret struct {
	V interface{}
	E error
}

func NewServer() *rpc.Server {
	mux := rpc.NewRespondMux()

	// most of this cruft can be eliminated when libhostbridge becomes thread-safe.
	// ideal: mux.Handle("window", fn.HandlerFrom(window.Module))
	mux.Handle("window.Create", fn.HandlerFrom(func(opts window.Options) (uintptr, error) {
		rchan := make(chan ret)
		app.Dispatch(func() {
			w, err := window.Create(opts)
			rchan <- ret{V: w, E: err}
		})
		r := <-rchan
		return uintptr(unsafe.Pointer(r.V.(*window.Window))), r.E
	}))

	return &rpc.Server{
		Codec:   codec.JSONCodec{},
		Handler: mux,
	}
}
