// client is the API client to bridge that you actually use
package client

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/qtalk-go/mux"
	"github.com/progrium/qtalk-go/talk"
	"github.com/progrium/qtalk-go/x/cbor/codec"
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
		!errors.Is(err, os.ErrClosed) &&
		!errors.Is(err, io.EOF) &&
		!errors.Is(err, syscall.EPIPE) &&
		!errors.Is(err, syscall.ECONNRESET) {
		return err
	}
	if err := c.Peer.Close(); err != nil &&
		!errors.Is(err, net.ErrClosed) &&
		!errors.Is(err, os.ErrClosed) &&
		!errors.Is(err, io.EOF) &&
		!errors.Is(err, syscall.EPIPE) &&
		!errors.Is(err, syscall.ECONNRESET) {
		return err
	}
	if c.cmd != nil {
		c.cmd.Process.Kill()
	}
	return nil
}

func (c *Client) Wait() error {
	return c.cmd.Wait()
}

func Bind(client *Client, name string, f reflect.Value) {
	f.Set(reflect.MakeFunc(f.Type(), func(args []reflect.Value) (result []reflect.Value) {
		ctx, _ := reflect.ValueOf(args[0]).Interface().(context.Context)
		_, err := client.Call(ctx, name, fn.Args{args[1]}, nil)
		return []reflect.Value{reflect.ValueOf(&err).Elem()}
	}))
}

func BindAll(client *Client, name string, p interface{}) {
	value := reflect.ValueOf(p).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)

		if field.Type.Kind() == reflect.Func && !strings.HasPrefix(field.Name, "On") {
			fullName := name + "." + field.Name
			log.Println(field.Index[0], field.Name, fullName)
			Bind(client, fullName, value.Field(i))
		}
	}
}

func New(peer *talk.Peer) *Client {
	client := &Client{Peer: peer}

	client.Window = &WindowModule{client: client, windows: make(map[Handle]*Window)}
	BindAll(client, "window", client.Window)
	client.System = &SystemModule{client: client}
	BindAll(client, "system", client.System)
	client.App = &AppModule{client: client}
	BindAll(client, "app", client.App)
	client.Menu = &MenuModule{client: client}
	BindAll(client, "menu", client.Menu)
	client.Shell = &ShellModule{client: client}
	BindAll(client, "shell", client.Shell)

	resp, err := client.Call(context.Background(), "Listen", nil, nil)
	if err == nil {
		go dispatchEvents(client, resp)
	}
	go client.Respond()
	return client
}

func Dial(addr string) (*Client, error) {
	peer, err := talk.Dial("tcp", addr, codec.CBORCodec{})
	if err != nil {
		return nil, err
	}
	return New(peer), nil
}

func findCmd() string {
	cmd := os.Getenv("APPTRON_CMD")
	if cmd == "" {
		if runtime.GOOS == "windows" {
			cmd = "apptron.exe"
		} else {
			cmd = "apptron"
		}
	}

	selfcmd, err := os.Executable()
	if err == nil && strings.Contains(strings.ToLower(selfcmd), cmd) {
		return selfcmd
	}

	dircmd := filepath.Join(".", cmd)
	if fi, err := os.Stat(dircmd); err == nil && fi.Mode().Perm()&0111 != 0 {
		return dircmd
	}

	pathcmd, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatal(err)
	}
	return pathcmd
}

func Spawn() (*Client, error) {
	cmd := exec.Command(findCmd(), "bridge")
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
	client := New(talk.NewPeer(sess, codec.CBORCodec{}))
	client.cmd = cmd
	return client, nil
}
