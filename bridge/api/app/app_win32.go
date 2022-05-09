package app

import (
	"fmt"
)

import (
	"tractor.dev/apptron/bridge/api/menu"
	"tractor.dev/apptron/bridge/win32"
	"tractor.dev/apptron/bridge/platform"
)

var (
	mainMenu *menu.Menu
)

func init() {
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
	fmt.Println("NewIndicator", icon)

  win32.SetProcessDpiAwarenessContext(win32.DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2)

	menu := menu.New(items)
	win32.SetupTray(menu.HMENU)

	/*
	m.Popup()
	*/

	for {
		win32.PollEvents()
	}

	/*
	obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
	obj.Retain()
	//obj.Button().SetTitle(i.Text)
	data := mac.NSData_WithBytes(icon, uint64(len(icon)))
	image := cocoa.NSImage_InitWithData(data)
	image.SetSize(mac.Size(16.0, 16.0))
	image.SetTemplate(true)
	obj.Button().SetImage(image)
	obj.Button().SetImagePosition(cocoa.NSImageOnly)

	menu := menu.New(items)
	obj.SetMenu(menu.NSMenu)
	*/
}

func Run(options Options) error {
  //win32.SetProcessDpiAwarenessContext(win32.DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2)
  /*
	win32.CreateTestWindow()

	for {
		win32.PollEvents()
	}
	*/

	platform.Start()
	return nil
}
