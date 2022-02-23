package client

import (
	"net"
	"runtime"
	"testing"

	"github.com/progrium/hostbridge/bridge"
	"github.com/progrium/hostbridge/bridge/core"
)

func init() {
	runtime.LockOSThread()
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	go func() {
		m.Run()
		core.Quit()
	}()
	core.Run(nil)
}

func setupBridgeClient(t *testing.T) (*Client, func()) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	srv := bridge.NewServer()
	go srv.Serve(l)

	client, err := Dial(l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	return client, func() {
		client.Close()
		l.Close()
	}
}
