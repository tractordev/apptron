package client

import (
	"context"

	"github.com/progrium/qtalk-go/fn"
)

type AppModule struct {
	client *Client
}

// Menu
func (m *AppModule) Menu(ctx context.Context) (ret []Display, err error) {
	_, err = m.client.Call(ctx, "app.Menu", fn.Args{}, &ret)
	return
}

// SetMenu
func (m *AppModule) SetMenu(ctx context.Context) (ret []Display, err error) {
	_, err = m.client.Call(ctx, "app.SetMenu", fn.Args{}, &ret)
	return
}

// NewIndicator
func (m *AppModule) NewIndicator(ctx context.Context) (ret []Display, err error) {
	_, err = m.client.Call(ctx, "app.NewIndicator", fn.Args{}, &ret)
	return
}
