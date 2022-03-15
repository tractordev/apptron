package client

import (
	"net"
	"runtime"
	"testing"

	"github.com/tractordev/hostbridge/bridge"
	"github.com/tractordev/hostbridge/bridge/core"
)

func init() {
	runtime.LockOSThread()
}

func TestMain(m *testing.M) {
	go func() {
		m.Run()
		core.Quit()
	}()
	core.Run(nil)
}

func setupBridgeClient(t *testing.T) (*Client, func()) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	srv := bridge.NewServer()
	go srv.Serve(l)

	client, err := Dial(l.Addr().String())
	if err != nil {
		panic(err)
	}

	return client, func() {
		client.Close()
		l.Close()
	}
}
