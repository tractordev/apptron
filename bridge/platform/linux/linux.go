//go:build linux

package linux

import (
	//"log"
	//"syscall"
	//"unsafe"
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

/*
type AppIndicator C.struct__AppIndicator
type GtkMenuShell C.struct__GtkMenuShell
type GtkWidget    C.struct__GtkWidget

func NewIndicator(id string, pngIconPath string, menu *GtkMenuShell) *AppIndicator {
	// @MemoryLeak: CString allocates memory but doesn't free?
	result := C.tray_indicator_new(C.CString(id), C.CString(pngIconPath), (*C.struct__GtkMenuShell)(menu))
	return (*AppIndicator)(result)
}

func MenuNew() *GtkMenuShell {
	result := C.menu_new()
	return (*GtkMenuShell)(result)
}

func MenuAppendMenuItem(menu *GtkMenuShell, item *GtkWidget) {
	C.menu_append_menu_item((*C.struct__GtkMenuShell)(menu), (*C.struct__GtkWidget)(item))
}

func MenuItemNew(id int, title string, disabled bool, checked bool, separator bool) *GtkWidget {
	result := C.menu_item_new(C.int(id), C.CString(title), toBool(disabled), toBool(checked), toBool(separator))
	return (*GtkWidget)(result)
}

func MenuItemSetSubmenu(parent *GtkWidget, child *GtkWidget) {
	C.menu_item_set_submenu((*C.struct__GtkWidget)(parent), (*C.struct__GtkWidget)(child))
}

func toBool(value bool) C.int {
	if (value) {
		return C.int(1)
	}
	return C.int(0)
}
*/

func TestNewIndicator() {
	/*
	C.tray_init()
	
	menu := C.menu_new()
	item := C.menu_item_new(1, C.CString("Hello, Sailor!"), C.int(0), C.int(0), C.int(0))
	C.menu_append_menu_item(menu, item)

	C.tray_indicator_new(C.CString("tray-id"), C.CString("/home/nick/dev/_projects/apptron/bridge/misc/icon.png"), menu)
	*/

	C.tray_test()
}
