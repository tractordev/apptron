// client is the API client to bridge that you actually use
package client

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/mux"
	"github.com/progrium/qtalk-go/rpc"
	"github.com/progrium/qtalk-go/talk"
)

type Client struct {
	*talk.Peer

	Window *WindowModule
	System *SystemModule
	Shell  *ShellModule
	App    *AppModule
	Menu   *MenuModule

	OnEvent func(event Event)

	files sync.Map
	cmd   *exec.Cmd
}

func (c *Client) Close() error {
	ctx := context.Background()
	if _, err := c.Call(ctx, "Shutdown", nil, nil); err != nil &&
		!errors.Is(err, net.ErrClosed) &&
		!errors.Is(err, io.EOF) &&
		!errors.Is(err, syscall.EPIPE) &&
		!errors.Is(err, syscall.ECONNRESET) {
		return err
	}
	if c.cmd != nil {
		c.cmd.Process.Kill()
	}
	if err := c.Peer.Close(); err != nil &&
		!errors.Is(err, net.ErrClosed) &&
		!errors.Is(err, io.EOF) &&
		!errors.Is(err, syscall.EPIPE) &&
		!errors.Is(err, syscall.ECONNRESET) {
		return err
	}
	return nil
}

func (c *Client) Wait() error {
	return c.cmd.Wait()
}

func (c *Client) ServeData(d []byte) string {
	hash := sha1.New()
	hash.Write(d)
	selector := hex.EncodeToString(hash.Sum(nil))
	dd, existed := c.files.LoadOrStore(selector, d)
	if !existed {
		c.Handle(selector, rpc.HandlerFunc(func(resp rpc.Responder, call *rpc.Call) {
			call.Receive(nil)
			ch, err := resp.Continue(nil)
			if err != nil {
				log.Println(err)
				return
			}
			defer ch.Close()
			buf := bytes.NewBuffer(dd.([]byte))
			if _, err := io.Copy(ch, buf); err != nil {
				log.Println(err)
				return
			}
		}))
	}
	return selector
}

func New(peer *talk.Peer) *Client {
	client := &Client{Peer: peer}
	client.Window = &WindowModule{client: client, windows: make(map[Handle]*Window)}
	client.System = &SystemModule{client: client}
	client.App = &AppModule{client: client}
	client.Menu = &MenuModule{client: client}
	client.Shell = &ShellModule{client: client}
	resp, err := client.Call(context.Background(), "Listen", nil, nil)
	if err == nil {
		go dispatchEvents(client, resp)
	}
	go client.Respond()
	return client
}

func Dial(addr string) (*Client, error) {
	peer, err := talk.Dial("tcp", addr, codec.JSONCodec{})
	if err != nil {
		return nil, err
	}
	return New(peer), nil
}

func Spawn() (*Client, error) {
	bridgecmd := os.Getenv("BRIDGECMD")
	if bridgecmd == "" {
		bridgecmd = "apptron"
	}
	cmdpath, err := exec.LookPath(bridgecmd)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(cmdpath)
	cmd.Stderr = os.Stderr
	wc, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	rc, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	sess, err := mux.DialIO(wc, rc)
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	client := New(talk.NewPeer(sess, codec.JSONCodec{}))
	client.cmd = cmd
	return client, nil
}
