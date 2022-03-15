package app

/*
#include "../../lib/hostbridge.h"
*/
import "C"

import (
	"context"
	"errors"
	"unsafe"

	"github.com/progrium/qtalk-go/rpc"
	"github.com/tractordev/hostbridge/bridge/core"
	"github.com/tractordev/hostbridge/bridge/menu"
)

var (
	Module       *module
	ErrBadHandle = errors.New("bad handle")
)

func init() {
	Module = &module{menu: menu.New(nil)}
}

type module struct {
	menu *menu.Menu
}

func Menu() *menu.Menu {
	return Module.Menu()
}

func (m module) Menu() *menu.Menu {
	return m.menu
}

func SetMenu(m *menu.Menu) {
	Module.menu = m
}

func (mod *module) SetMenu(menuID core.Handle) error {
	var m = menu.Module.FindByID(menuID)

	if m == nil {
		return ErrBadHandle
	}

	mod.menu = m
	return nil
}

func NewIndicator(icon []byte, items []menu.Item) {
	var cicon C.Icon
	if len(icon) > 0 {
		cicon = C.Icon{data: (*C.uchar)(unsafe.Pointer(&icon[0])), size: C.int(len(icon))}
	} else {
		cicon = C.Icon{data: (*C.uchar)(nil), size: C.int(0)}
	}

	trayMenu := NewContextMenu(items)

	eventLoop := *(*C.EventLoop)(core.EventLoop())
	C.tray_set_system_tray(eventLoop, cicon, trayMenu)
}

func (m module) NewIndicator(iconSel string, items []menu.Item, call *rpc.Call) {
	var icon []byte
	icon, _ = core.FetchData(context.Background(), call, iconSel)
	core.Dispatch(func() {
		NewIndicator(icon, items)
	})
}

func NewContextMenu(items []menu.Item) C.ContextMenu {
	result := C.context_menu_create()

	for _, it := range items {
		if len(it.SubMenu) > 0 {
			submenu := NewContextMenu(it.SubMenu)
			C.context_menu_add_submenu(result, C.CString(it.Title), toCBool(it.Enabled), submenu)
		} else {
			C.context_menu_add_item(result, buildCMenuItem(it))
		}
	}

	return result
}

func buildCMenuItem(item menu.Item) C.Menu_Item {
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
