package app

import (
	"github.com/progrium/macdriver/cocoa"
	mac "github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"tractor.dev/apptron/bridge/api/menu"
	"tractor.dev/apptron/bridge/event"
	"tractor.dev/apptron/bridge/platform"
)

var (
	mainMenu *menu.Menu
	app      cocoa.NSApplication
)

func init() {
	app = cocoa.NSApp()
}

func Menu() *menu.Menu {
	return mainMenu
}

func SetMenu(menu *menu.Menu) error {
	app.SetMainMenu(menu.NSMenu)
	mainMenu = menu
	return nil
}

func NewIndicator(icon []byte, items []menu.Item) {
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
}

func Run(options Options) error {
	mainBundle := cocoa.NSBundle_Main()
	bundleClass := mainBundle.Class()
	bundleClass.AddMethod("__bundleIdentifier", func(self objc.Object) objc.Object {
		if self.Pointer() == mainBundle.Pointer() {
			return mac.String(options.Identifier)
		}
		// After the swizzle this will point to the original method, and return the
		// original bundle identifier.
		return self.Send("__bundleIdentifier")
	})
	bundleClass.Swizzle("bundleIdentifier", "__bundleIdentifier")

	DelegateClass := objc.NewClass("AppDelegate", "NSObject")
	DelegateClass.AddMethod("applicationShouldTerminateAfterLastWindowClosed:", func(notification objc.Object) bool {
		return !options.RunsAfterLastWindow
	})
	DelegateClass.AddMethod("applicationWillFinishLaunching:", func(notification objc.Object) {
		if mainMenu == nil {
			mainMenu = menu.New([]menu.Item{})
		}
		app.SetMainMenu(mainMenu.NSMenu)
		if options.AccessoryMode {
			app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyAccessory)
		} else {
			app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyRegular)
		}
	})
	DelegateClass.AddMethod("applicationDidFinishLaunching:", func(notification objc.Object) {
		app.ActivateIgnoringOtherApps(true)
	})
	DelegateClass.AddMethod("menuClick:", func(self, sender objc.Object) {
		event.Emit(event.Event{
			Type:     event.MenuItem,
			MenuItem: int(sender.Get("tag").Int()),
		})
	})
	objc.RegisterClass(DelegateClass)

	app.SetDelegate(objc.Get("AppDelegate").Alloc().Init())

	platform.Start()
	return nil
}
