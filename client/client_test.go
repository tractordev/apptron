package client

import (
	"net"
	"runtime"
	"testing"

	"tractor.dev/apptron/bridge"
	"tractor.dev/apptron/bridge/platform"
)

func init() {
	runtime.LockOSThread()
}

func TestMain(m *testing.M) {
	go func() {
		m.Run()
		platform.Terminate()
	}()
	platform.Main()
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
