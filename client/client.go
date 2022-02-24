package client

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"sync"

	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/rpc"
	"github.com/progrium/qtalk-go/talk"
)

type Client struct {
	*talk.Peer

	Window *WindowModule
	Screen *ScreenModule

	files sync.Map
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

func Dial(addr string) (*Client, error) {
	peer, err := talk.Dial("tcp", addr, codec.JSONCodec{})
	if err != nil {
		return nil, err
	}
	client := &Client{Peer: peer}
	client.Window = &WindowModule{client: client}
	client.Screen = &ScreenModule{client: client}
	resp, err := client.Call(context.Background(), "Listen", nil, nil)
	if err != nil {
		return nil, err
	}
	go dispatchEvents(client, resp)
	return client, nil
}
