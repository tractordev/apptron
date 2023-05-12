package menu

import (
	"tractor.dev/apptron/bridge/misc"
	"tractor.dev/apptron/bridge/resource"
)

var Module *module

type module struct {
	OnClick func(event Event)
}

func init() {
	Module = &module{}
}

var (
	mainMenu *Menu
)

func Main() *Menu {
	return mainMenu
}

func SetMain(menu *Menu) error {
	mainMenu = menu
	return nil
}

func Get(handle resource.Handle) (*Menu, error) {
	v, err := resource.Lookup(handle)
	if err != nil {
		return nil, err
	}
	w, ok := v.(*Menu)
	if !ok {
		return nil, resource.ErrBadHandle
	}
	return w, nil
}

type menu struct {
	Handle resource.Handle
	Items  []Item
}

type Item struct {
	ID          int
	Title       string
	Disabled    bool
	Selected    bool
	Separator   bool
	Accelerator string
	SubMenu     []Item
}

type Position = misc.Position

func (m *module) New(items []Item) *Menu {
	mm := New(items)
	resource.Retain(mm.Handle, mm)
	return mm
}

func (mm *module) Destroy(h resource.Handle) (err error) {
	var m *Menu
	if m, err = Get(h); err == nil {
		m.Destroy()
		resource.Release(h)
	}
	return
}

func (mm *module) Popup(items []Item) int {
	m := New(items)
	defer m.Destroy()
	return m.Popup()
}
