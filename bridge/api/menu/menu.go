package menu

import (
	"tractor.dev/hostbridge/bridge/resource"
)

var Module *module

type module struct{}

func init() {
	Module = &module{}
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

func (m *module) New(items []Item) *Menu {
	return New(items)
}
