//go:build linux

package linux

import (
	//"log"
	//"syscall"
	"unsafe"
)

/*
#cgo linux pkg-config: appindicator3-0.1

#include "linux.h"
*/
import "C"

func OS_Init() {
	C.tray_init()
}

func PollEvents() {
	C.tray_poll_events()
}

type MenuHandle uintptr

func NewIndicator(id string, pngIconPath string, menu MenuHandle) *C.struct__AppIndicator {
	// @MemoryLeak: CString allocates memory but doesn't free?
	result := C.tray_indicator_new(C.CString(id), C.CString(pngIconPath), fromHandle(menu))
	return result
}

func MenuNew() MenuHandle {
	//C.tray_init()

	result := C.menu_new()
	return toHandle(result)
}

func MenuAppendMenuItem(menu MenuHandle, item *C.struct__GtkWidget) {
	C.menu_append_menu_item(fromHandle(menu), item)
}

func MenuItemNew(id int, title string, disabled bool, checked bool, separator bool) *C.struct__GtkWidget {
	result := C.menu_item_new(C.int(id), C.CString(title), toBool(disabled), toBool(checked), toBool(separator))
	return result
}

func MenuItemSetSubmenu(parent *C.struct__GtkWidget, child MenuHandle) {
	C.menu_item_set_submenu(parent, fromHandle(child))
}

func toBool(value bool) C.int {
	if (value) {
		return C.int(1)
	}
	return C.int(0)
}

func toHandle(menu *C.struct__GtkMenuShell) MenuHandle {
	return (MenuHandle)(unsafe.Pointer(menu))
}

func fromHandle(menu MenuHandle) *C.struct__GtkMenuShell {
	return (*C.struct__GtkMenuShell)(unsafe.Pointer(menu))
}
