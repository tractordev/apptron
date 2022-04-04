package client

import (
	"context"

	"github.com/progrium/qtalk-go/fn"
)

type AppModule struct {
	client *Client
}

// Run
func (m *AppModule) Run(ctx context.Context) (err error) {
	_, err = m.client.Call(ctx, "app.Run", fn.Args{}, nil)
	return
}

// Menu
func (m *AppModule) Menu(ctx context.Context) (ret Menu, err error) {
	_, err = m.client.Call(ctx, "app.Menu", fn.Args{}, &ret)
	return
}

// SetMenu
func (m *AppModule) SetMenu(ctx context.Context, menu Menu) (err error) {
	_, err = m.client.Call(ctx, "app.SetMenu", fn.Args{menu.Handle}, nil)
	return
}

// NewIndicator
func (m *AppModule) NewIndicator(ctx context.Context, icon []byte, items []MenuItem) (err error) {
	iconSel := m.client.ServeData(icon)
	_, err = m.client.Call(ctx, "app.NewIndicator", fn.Args{iconSel, items}, nil)
	return
}
