package bridge

import (
	"unsafe"

	"github.com/mitchellh/mapstructure"
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

	// most of this cruft can be eliminated when qtalk gets mapstructure support
	// and libhostbridge becomes thread-safe.
	mux.Handle("window.Create", fn.HandlerFrom(func(opts map[string]interface{}) (uintptr, error) {
		var wopts window.Options
		err := mapstructure.Decode(opts, &wopts)
		if err != nil {
			return 0, err
		}
		rchan := make(chan ret)
		app.Dispatch(func() {
			w, err := window.Create(wopts)
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
