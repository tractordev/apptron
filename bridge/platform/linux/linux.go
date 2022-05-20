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
	C.gtk_init_check(nil, nil)
}

func PollEvents() {
	C.gtk_main_iteration_do(0)
}

func NewIndicator(id string, pngIconPath string, menu MenuHandle) IndicatorHandle {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	result := C.app_indicator_new(cid, C.CString(""), C.APP_INDICATOR_CATEGORY_APPLICATION_STATUS)
	C.app_indicator_set_status(result, C.APP_INDICATOR_STATUS_ACTIVE)

	//app_indicator_set_title(global_app_indicator, title);
	//app_indicator_set_label(global_app_indicator, title, "");

	if len(pngIconPath) > 0 {
		cIconPath := C.CString(pngIconPath)
		defer C.free(unsafe.Pointer(cIconPath))

		C.app_indicator_set_icon_full(result, cIconPath, C.CString(""))
	}

	if menu != 0 {
		C.app_indicator_set_menu(result, (*C.struct__GtkMenu)(unsafe.Pointer(fromMenuHandle(menu))))
	}

	return toIndicatorHandle(result)
}

func MenuNew() MenuHandle {
	result := C.gtk_menu_new()
	return toMenuHandle(result)
}

func MenuAppendMenuItem(menu MenuHandle, item MenuItemHandle) {
	menuShell := (*C.struct__GtkMenuShell)(unsafe.Pointer(fromMenuHandle(menu)))
	C.gtk_menu_shell_append(menuShell, fromMenuItemHandle(item))
}

func MenuItemNew(id int, title string, disabled bool, checked bool, separator bool) MenuItemHandle {
	var result *C.struct__GtkWidget = nil

	if separator {
		result = C.gtk_separator_menu_item_new()
		C.gtk_widget_show(result)
	} else {
		ctitle := C.CString(title)
		defer C.free(unsafe.Pointer(ctitle))

		if checked {
			result = C.gtk_check_menu_item_new_with_label(ctitle)

			checkMenuItem := (*C.struct__GtkCheckMenuItem)(unsafe.Pointer(result))
			C.gtk_check_menu_item_set_active(checkMenuItem, toBool(checked))
		} else {
			result = C.gtk_menu_item_new_with_label(ctitle)
		}

		C.gtk_widget_set_sensitive(result, toBool(!disabled))

	    //
	    // NOTE(nick): accelerators seem to require a window and an accel_group
	    // Are they even supported in the AppIndicator?
	    // As far as I can tell they don't ever show up in the AppIndicator menu...
	    //
	    // @see https://github.com/bstpierre/gtk-examples/blob/master/c/accel.c
	    //
	    /*
	    GtkWindow *window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
	    GtkAccelGroup *accel_group = gtk_accel_group_new();
	    gtk_window_add_accel_group(GTK_WINDOW(window), accel_group);

	    gtk_widget_add_accelerator(item, "activate", accel_group, GDK_KEY_F7, 0, GTK_ACCEL_VISIBLE);
	    */

	    cactivate := C.CString("activate")
	    defer C.free(unsafe.Pointer(cactivate))

	    C._g_signal_connect(result, cactivate, C.go_menu_callback, C.int(id))

	    C.gtk_widget_show(result)
	}

	return toMenuItemHandle(result)
}

func MenuItemSetSubmenu(parent MenuItemHandle, child MenuHandle) {
	menuItem := (*C.struct__GtkMenuItem)(unsafe.Pointer(fromMenuItemHandle(parent)))
	C.gtk_menu_item_set_submenu(menuItem, fromMenuHandle(child));
}

//export go_menu_callback
func go_menu_callback(item *C.struct__GtkMenuItem, menuId C.int) {
    if globalMenuCallback != nil {
    	globalMenuCallback(int(menuId))
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

func toMenuHandle(menu *C.struct__GtkWidget) MenuHandle {
	return (MenuHandle)(unsafe.Pointer(menu))
}

func fromMenuHandle(menu MenuHandle) *C.struct__GtkWidget {
	return (*C.struct__GtkWidget)(unsafe.Pointer(menu))
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
