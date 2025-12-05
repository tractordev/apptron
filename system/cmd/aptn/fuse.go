package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/hugelgupf/p9/p9"
	"tractor.dev/toolkit-go/engine/cli"
	"tractor.dev/wanix/fs/fusekit"
	"tractor.dev/wanix/fs/p9kit"
	"tractor.dev/wanix/vm/v86/shm"
)

func fuseCmd() *cli.Command {
	return &cli.Command{
		Usage: "fuse",
		Short: "mount experimental fuse filesystem",
		Run:   setupFuseFS,
	}
}

func setupFuseFS(ctx *cli.Context, args []string) {
	sch, err := shm.NewSharedChannel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create channel: %v\n", err)
		os.Exit(1)
	}
	defer sch.Close()
	fsys, err := p9kit.ClientFS(&rwcConn{rwc: sch}, "/", p9.WithMessageSize(512*1024))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create client FS: %v\n", err)
		os.Exit(1)
	}
	os.MkdirAll("/x", 0755)
	mount, err := fusekit.Mount(fsys, "/x", context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mount.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)
	for sig := range sigChan {
		if sig == os.Interrupt {
			return
		}
	}

	select {}
}

// Conn is an adapter that implements net.Conn using an underlying io.ReadWriteCloser.
// LocalAddr/RemoteAddr will be dummy addrs, SetDeadline/Set[Read|Write]Deadline are no-ops.
type rwcConn struct {
	rwc io.ReadWriteCloser
}

func (c *rwcConn) Read(b []byte) (int, error) {
	return c.rwc.Read(b)
}
func (c *rwcConn) Write(b []byte) (int, error) {
	return c.rwc.Write(b)
}
func (c *rwcConn) Close() error {
	return c.rwc.Close()
}
func (c *rwcConn) LocalAddr() (addr net.Addr) {
	return dummyAddr("rwc-local")
}
func (c *rwcConn) RemoteAddr() (addr net.Addr) {
	return dummyAddr("rwc-remote")
}
func (c *rwcConn) SetDeadline(t time.Time) error {
	return nil // not supported
}
func (c *rwcConn) SetReadDeadline(t time.Time) error {
	return nil // not supported
}
func (c *rwcConn) SetWriteDeadline(t time.Time) error {
	return nil // not supported
}

type dummyAddr string

func (a dummyAddr) Network() string { return string(a) }
func (a dummyAddr) String() string  { return string(a) }
