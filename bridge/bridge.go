package bridge

import (
	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/qtalk-go/rpc"

	"github.com/progrium/hostbridge/bridge/window"
)

func NewServer() *rpc.Server {
	mux := rpc.NewRespondMux()

	mux.Handle("window", fn.HandlerFrom(&window.Module{}))

	return &rpc.Server{
		Codec:   codec.JSONCodec{},
		Handler: mux,
	}
}
