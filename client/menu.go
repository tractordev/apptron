package client

import (
	"context"

	"github.com/progrium/qtalk-go/fn"
)

type Menu struct {
	/*
		Items []Item
	*/
}

type Item struct {
	ID          uint16
	Title       string
	Enabled     bool
	Selected    bool
	Accelerator string
	SubMenu     []Item
}

type MenuModule struct {
	client *Client
}

// New
func (m *MenuModule) New(ctx context.Context, items []Item) (ret interface{}, err error) {
	_, err = m.client.Call(ctx, "menu.New", fn.Args{items}, &ret)
	return
}
