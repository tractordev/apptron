package app

import (
	"context"

	"github.com/progrium/qtalk-go/rpc"
	"tractor.dev/apptron/bridge/api/menu"
	"tractor.dev/apptron/bridge/misc"
	"tractor.dev/apptron/bridge/resource"
)

var Module *module

type module struct{}

func init() {
	Module = &module{}
}

type Options struct {
	Identifier string
	Agent      bool
	Accessory  bool
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

func (m *module) NewIndicator(iconSel string, items []menu.Item, call *rpc.Call) error {
	icon, err := misc.FetchData(context.Background(), call, iconSel)
	if err != nil {
		return err
	}
	NewIndicator(icon, items)
	return nil
}

func (m *module) Run(options Options) error {
	if options.Identifier == "" {
		options.Identifier = "com.progrium.Apptron"
	}
	return Run(options)
}
