package menu

/*
#include "../../lib/hostbridge.h"
*/
import "C"

import (
	"sync"

	"github.com/progrium/hostbridge/bridge/core"
)

var Module *module

func init() {
	Module = &module{}
}

type module struct {
	mu sync.Mutex

	menus      []Menu
	nextMenuId int 
}

type Menu struct {
	ID     core.Handle
	Handle C.Menu

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

	/*
		Role        string // for wry's add_native_item (see Electron's MenuItem role for examples)
		Type        string // normal, separator, submenu, checkbox or radio
	*/

	SubMenu []Item
}

func FindByID(menuID core.Handle) *Menu {
	return Module.FindByID(menuID)
}

func (m *module) FindByID(menuID core.Handle) *Menu {
	m.mu.Lock()
	defer m.mu.Unlock()

	var index int = -1

	for i, v := range m.menus {
		if v.ID == menuID {
			index = i
			break
		}
	}

	if index >= 0 {
		return &m.menus[index]
	}

	return nil
}

func New(items []Item) *Menu {
	return Module.New(items)
}

func (m module) New(items []Item) *Menu {
	cmenu := C.menu_create()

	for _, it := range items {
		if len(it.SubMenu) > 0 {
			submenu := m.New(it.SubMenu)
			C.menu_add_submenu(cmenu, C.CString(it.Title), toCBool(it.Enabled), submenu.Handle)
		} else {
			C.menu_add_item(cmenu, buildCMenuItem(it))
		}
	}

	var id = -1

	m.mu.Lock()
	m.nextMenuId += 1
	id = m.nextMenuId
	m.mu.Unlock()

	menu := &Menu{}
	menu.Handle = cmenu
	menu.ID     = core.Handle(id)

	return menu
}

func buildCMenuItem(item Item) C.Menu_Item {
	return C.Menu_Item{
		id:          C.int(item.ID),
		title:       C.CString(item.Title),
		enabled:     toCBool(item.Enabled),
		selected:    toCBool(item.Selected),
		accelerator: C.CString(item.Accelerator),
	}
}

func toCBool(it bool) C.uchar {
	if it {
		return C.uchar(1)
	}

	return C.uchar(0)
}
