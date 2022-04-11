package menu

import (
	"io"
	"time"

	"github.com/progrium/macdriver/cocoa"
	mac "github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"tractor.dev/apptron/bridge/event"
	"tractor.dev/apptron/bridge/resource"
)

type Menu struct {
	menu
	cocoa.NSMenu `json:"-"`
}

func New(items []Item) *Menu {
	menu := &Menu{
		menu: menu{
			Handle: resource.NewHandle(),
			Items:  items,
		},
		NSMenu: cocoa.NSMenu_New(),
	}

	menu.SetAutoenablesItems(true)

	for _, i := range items {
		menu.AddItem(newMenuItem(i))
	}

	return menu
}

func (m *Menu) Destroy() {
	m.NSMenu.Release()
}

func (m *Menu) Popup() int {
	ch := make(chan int, 1)
	event.Listen(time.Now(), func(e event.Event) error {
		if e.Type == event.MenuItem {
			ch <- e.MenuItem
		}
		return io.EOF
	})
	if m.NSMenu.PopUpMenuPositioningItem_atLocation_inView_(nil, cocoa.NSEvent_mouseLocation(), nil) {
		return <-ch
	}
	return 0
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

	if !i.Disabled && len(i.SubMenu) == 0 {
		// found in AppDelegate
		item.SetAction(objc.Sel("menuClick:"))

		// special item titles
		if i.Title == "Quit" {
			item.SetTarget(cocoa.NSApp())
			item.SetAction(objc.Sel("terminate:"))
		}
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
