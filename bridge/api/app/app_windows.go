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
	// @see https://github.com/glfw/glfw/blob/master/src/win32_init.c#L692

	// NOTE(nick): the exact snippet from GLFW is:
	/*
			if (_glfwIsWindows10Version1703OrGreaterWin32())
		        SetProcessDpiAwarenessContext(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2);
		    else if (IsWindows8Point1OrGreater())
		        SetProcessDpiAwareness(PROCESS_PER_MONITOR_DPI_AWARE);
		    else if (IsWindowsVistaOrGreater())
		        SetProcessDPIAware();

	*/
	// BUT, I think it's sufficient to just check if these proceedures are loaded?
	// @Robustness: test this assumption

	success := win32.SetProcessDpiAwarenessContext(win32.DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2)
	if !success {
		success = win32.SetProcessDpiAwareness(win32.PROCESS_PER_MONITOR_DPI_AWARE)

		if !success {
			success = win32.SetProcessDPIAware()
		}
	}
}

func Menu() *menu.Menu {
	return mainMenu
}

func SetMenu(menu *menu.Menu) error {
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

	win32.NewTrayMenu(menu.PopupMenu, icon, onClick)
}

func Run(options Options) error {
	platform.Start()
	return nil
}
