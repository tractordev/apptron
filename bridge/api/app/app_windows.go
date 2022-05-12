package app

import (
	"tractor.dev/apptron/bridge/api/menu"
	"tractor.dev/apptron/bridge/event"
	"tractor.dev/apptron/bridge/platform"
	"tractor.dev/apptron/bridge/platform/win32"
)

var (
	mainMenu *menu.Menu
)

func init() {
	//
	// @Robustness: add support for older versions of Windows
	// @see https://github.com/glfw/glfw/blob/master/src/win32_init.c#L643
	//
	win32.SetProcessDpiAwarenessContext(win32.DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2)
}

func Menu() *menu.Menu {
	return mainMenu
}

func SetMenu(menu *menu.Menu) error {
	//app.SetMainMenu(menu.NSMenu)
	mainMenu = menu
	return nil
}

func NewIndicator(icon []byte, items []menu.Item) {
	menu := menu.New(items)
	onClick := func(id int32) {
		event.Emit(event.Event{
			Type:     event.MenuItem,
			MenuItem: int(id),
		})
	}
	win32.SetTrayMenu(menu.HMENU, icon, onClick)
}

func Run(options Options) error {
	platform.Start()
	return nil
}
