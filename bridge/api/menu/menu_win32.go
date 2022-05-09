package menu

import (
	"tractor.dev/apptron/bridge/resource"
)

type Menu struct {
	menu
}

func New(items []Item) *Menu {
	menu := &Menu{
		menu: menu{
			Handle: resource.NewHandle(),
			Items:  items,
		},
	}

	return menu
}

func (m *Menu) Destroy() {
}

func (m *Menu) Popup() int {
	return 0
}
