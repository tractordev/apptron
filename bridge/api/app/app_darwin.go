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
	app cocoa.NSApplication
)

func init() {
	app = cocoa.NSApp()
}

func SetMenu(m *menu.Menu) error {
	app.SetMainMenu(m.NSMenu)
	return menu.SetMenu(m)
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
		return !options.Agent
	})
	DelegateClass.AddMethod("applicationWillFinishLaunching:", func(notification objc.Object) {
		mainMenu := menu.GetMenu()
		if mainMenu == nil {
			menu.SetMenu(menu.New([]menu.Item{}))
			mainMenu = menu.GetMenu()
			app.SetMainMenu(mainMenu.NSMenu)
		}
		if options.Accessory {
			app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyAccessory)
		} else {
			app.SetMainMenu(mainMenu.NSMenu)
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

	if options.DisableAutoSave != true {
		setupWindowRestoreListener(options.Identifier)
	}

	platform.Start()
	return nil
}
