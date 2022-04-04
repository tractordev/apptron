package app

import (
	"context"

	"github.com/progrium/qtalk-go/rpc"
	"tractor.dev/hostbridge/bridge/api/menu"
	"tractor.dev/hostbridge/bridge/misc"
	"tractor.dev/hostbridge/bridge/resource"
)

var Module *module

type module struct{}

func init() {
	Module = &module{}
}

func (m *module) Menu() *menu.Menu {
	return Menu()
}

func (m *module) SetMenu(handle resource.Handle) error {
	mm, err := menu.Get(handle)
	if err != nil {
		return err
	}
	SetMenu(mm)
	return nil
}

func (m *module) NewIndicator(iconSel string, items []menu.Item, call *rpc.Call) {
	var icon []byte
	icon, _ = misc.FetchData(context.Background(), call, iconSel)
	NewIndicator(icon, items)
}

func (m *module) Run() error {
	return Run()
}
