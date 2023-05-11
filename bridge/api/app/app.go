package app

import (
	"tractor.dev/apptron/bridge/api/menu"
)

var Module *module

type module struct{}

func init() {
	Module = &module{}
}

type Options struct {
	Identifier      string
	Agent           bool // app should not terminate when last window closes
	Accessory       bool // app should not be task switchable
	DisableAutoSave bool // disable window position saving and restoring
}

func (m *module) Menu() *menu.Menu {
	return menu.Main()
}

func (m *module) SetMenu(men *menu.Menu) {
	menu.SetMain(men)
}

func (m *module) NewIndicator(icon []byte, items []menu.Item) error {
	NewIndicator(icon, items)
	return nil
}

func (m *module) Run(options Options) error {
	if options.Identifier == "" {
		options.Identifier = "com.progrium.Apptron"
	}
	return Run(options)
}
