package menu

/*
#include "../../lib/hostbridge.h"
*/
import "C"

var Module *module

func init() {
	Module = &module{}
}

type module struct{}

type Menu struct {
	ID core.Handle // Library handle, not user id.

	/*
		Items []Item
	*/

	Handle C.Menu
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

func New(items []Item) *Menu {
	return Module.New(items)
}

func (m module) New(items []Item) *Menu {
	cmenu_id := C.menu_create()

	for _, it := range items {
		if len(it.SubMenu) > 0 {
			submenu := m.New(it.SubMenu)
			C.menu_add_submenu(cmenu_id, C.CString(it.Title), toCBool(it.Enabled), C.int(submenu.ID))
		} else {
			C.menu_add_item(cmenu_id, buildCMenuItem(it))
		}
	}

	menu := &Menu{}
	menu.ID = core.Handle(cmenu_id)

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
