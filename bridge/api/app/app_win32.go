package app

import (
	"tractor.dev/apptron/bridge/api/menu"
	//"tractor.dev/apptron/bridge/event"
	//"tractor.dev/apptron/bridge/platform"
)

var (
	mainMenu *menu.Menu
	//app      cocoa.NSApplication
)

func init() {
	//app = cocoa.NSApp()
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
	return nil
}
