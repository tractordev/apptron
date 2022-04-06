package client

import (
	"context"

	"github.com/progrium/qtalk-go/fn"
)

type Menu struct {
	Handle Handle

	/*
		Items []Item
	*/
}

type MenuItem struct {
	ID          int
	Title       string
	Disabled    bool
	Selected    bool
	Separator   bool
	Accelerator string
	SubMenu     []MenuItem
}

type MenuModule struct {
	client *Client

	OnClick func(event Event)
}

// New
func (m *MenuModule) New(ctx context.Context, items []MenuItem) (ret Menu, err error) {
	_, err = m.client.Call(ctx, "menu.New", fn.Args{items}, &ret)
	return
}

// Popup
func (m *MenuModule) Popup(ctx context.Context, menu Menu) (ret int, err error) {
	_, err = m.client.Call(ctx, "menu.Popup", fn.Args{menu.Handle}, &ret)
	return
}
