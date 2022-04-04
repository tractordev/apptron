package menu

import (
	"fmt"
	"log"
	"os"

	"github.com/progrium/macdriver/cocoa"
	mac "github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"tractor.dev/hostbridge/bridge/resource"
)

type Menu struct {
	menu
	cocoa.NSMenu `json:"-"`
}

// WIP
func (m *module) Popup(handle resource.Handle) {
	menu, err := Get(handle)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: use mouse cursor position
	// BLOCKS! handle somehow?
	fmt.Fprintln(os.Stderr, menu.NSMenu.PopUpMenuPositioningItem_atLocation_inView_(nil, mac.Point(100, 100), nil))
}

func New(items []Item) *Menu {
	menu := &Menu{
		menu: menu{
			Handle: resource.NewHandle(),
			Items:  items,
		},
		NSMenu: cocoa.NSMenu_New(),
	}
	resource.Retain(menu.Handle, menu)

	menu.SetAutoenablesItems(true)

	for _, i := range items {
		menu.AddItem(newMenuItem(i))
	}

	return menu
}

func newMenuItem(i Item) cocoa.NSMenuItem {
	if i.Separator {
		return cocoa.NSMenuItem_Separator()
	}

	item := cocoa.NSMenuItem_New()
	item.SetTitle(i.Title)
	item.SetTag_(mac.NSInteger(int(i.ID)))
	item.SetEnabled(!i.Disabled)
	// item.SetToolTip(i.Tooltip)

	// Checked
	if i.Selected {
		item.SetState(cocoa.NSControlStateValueOn)
	}

	// Icon
	// if i.Icon != "" {
	// 	b, err := base64.StdEncoding.DecodeString(i.Icon)
	// 	if err == nil {
	// 		data := core.NSData_WithBytes(b, uint64(len(b)))
	// 		img := cocoa.NSImage_InitWithData(data)
	// 		img.SetSize(core.Size(16, 16))
	// 		item.SetImage(img)
	// 	}
	// }

	// found in AppDelegate
	item.SetAction(objc.Sel("menuClick:"))

	// special item titles
	if i.Title == "Quit" {
		item.SetTarget(cocoa.NSApp())
		item.SetAction(objc.Sel("terminate:"))
	}

	if len(i.SubMenu) > 0 {
		sub := cocoa.NSMenu_New()
		sub.SetTitle(i.Title)
		sub.SetAutoenablesItems(true)
		for _, i := range i.SubMenu {
			sub.AddItem(newMenuItem(i))
		}
		item.SetSubmenu(sub)
	}

	return item
}
