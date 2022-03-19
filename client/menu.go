package client

import (
	"context"

	"github.com/progrium/qtalk-go/fn"
)

type Menu struct {
	ID Handle

	/*
		Items []Item
	*/
}

type MenuItem struct {
	ID          uint16
	Title       string
	Enabled     bool
	Selected    bool
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
