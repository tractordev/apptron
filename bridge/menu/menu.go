package menu

/*
#include "../../lib/hostbridge.h"
*/
import "C"

type Handle int

type Menu struct {
	/*
	Items []Item
	*/

	Handle C.Menu
}

type TrayMenu struct {
	//Handle C.Tray
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

	SubMenu     []Item
}

var AppMenu       Menu
var AppMenuWasSet bool

func init() {
	AppMenu = New([]Item {})
	AppMenuWasSet = false
}

func New(items []Item) Menu {
	menu := C.menu_create()

	for _, it := range items {
		if (len(it.SubMenu) > 0) {
			submenu := New(it.SubMenu)
			C.menu_add_submenu(menu, C.CString(it.Title), toCBool(it.Enabled), submenu.Handle)
		} else {
			C.menu_add_item(menu, buildCMenuItem(it))
		}
	}

	result := Menu{}
	result.Handle = menu

	return result
}

func buildCMenuItem(item Item) C.Menu_Item {
	return C.Menu_Item {
		id:          C.int(item.ID),
		title:       C.CString(item.Title),
		enabled:     toCBool(item.Enabled),
		selected:    toCBool(item.Selected),
		accelerator: C.CString(item.Accelerator),
	}
}

func toCBool(it bool) C.uchar {
	if (it) {
		return C.uchar(1)
	}

	return C.uchar(0)
}

func toBool(it C.uchar) bool {
	return int(it) != 0
}