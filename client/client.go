package client

import (
	"context"

	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/talk"
)

type Client struct {
	*talk.Peer

	Window *WindowModule
	Screen *ScreenModule
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
