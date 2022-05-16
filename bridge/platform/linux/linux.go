//go:build linux

package linux

import (
	"unsafe"
)

/*
#cgo linux pkg-config: appindicator3-0.1

#include "linux.h"
*/
import "C"

type Menu_Callback func(menuId int)

var globalMenuCallback Menu_Callback

//
// Exports
//

func OS_Init() {
	C.tray_init()
}

func PollEvents() {
	C.tray_poll_events()
}

func NewIndicator(id string, pngIconPath string, menu MenuHandle) IndicatorHandle {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	cIconPath := C.CString(pngIconPath)
	defer C.free(unsafe.Pointer(cIconPath))

	result := C.tray_indicator_new(cid, cIconPath, fromMenuHandle(menu))

	return toIndicatorHandle(result)
}

func MenuNew() MenuHandle {
	result := C.menu_new()
	return toMenuHandle(result)
}

func MenuAppendMenuItem(menu MenuHandle, item MenuItemHandle) {
	C.menu_append_menu_item(fromMenuHandle(menu), fromMenuItemHandle(item))
}

func MenuItemNew(id int, title string, disabled bool, checked bool, separator bool) MenuItemHandle {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	result := C.menu_item_new(C.int(id), ctitle, toBool(disabled), toBool(checked), toBool(separator))
	return toMenuItemHandle(result)
}

func MenuItemSetSubmenu(parent MenuItemHandle, child MenuHandle) {
	C.menu_item_set_submenu(fromMenuItemHandle(parent), fromMenuHandle(child))
}

//export go_menu_callback
func go_menu_callback(menuId int) {
    if globalMenuCallback != nil {
    	globalMenuCallback(menuId)
    }
}

func SetGlobalMenuCallback(callback Menu_Callback) {
	globalMenuCallback = callback
}

//
// Helpers
//

func toBool(value bool) C.int {
	if (value) {
		return C.int(1)
	}
	return C.int(0)
}

func toMenuHandle(menu *C.struct__GtkMenuShell) MenuHandle {
	return (MenuHandle)(unsafe.Pointer(menu))
}

func fromMenuHandle(menu MenuHandle) *C.struct__GtkMenuShell {
	return (*C.struct__GtkMenuShell)(unsafe.Pointer(menu))
}

func toMenuItemHandle(item *C.struct__GtkWidget) MenuItemHandle {
	return (MenuItemHandle)(unsafe.Pointer(item))
}

func fromMenuItemHandle(item MenuItemHandle) *C.struct__GtkWidget {
	return (*C.struct__GtkWidget)(unsafe.Pointer(item))
}

func toIndicatorHandle(indicator *C.struct__AppIndicator) IndicatorHandle {
	return (IndicatorHandle)(unsafe.Pointer(indicator))
}

func fromIndicatorHandle(indicator IndicatorHandle) *C.struct__AppIndicator {
	return (*C.struct__AppIndicator)(unsafe.Pointer(indicator))
}
