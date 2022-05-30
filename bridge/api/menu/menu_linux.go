package menu

import (
  "tractor.dev/apptron/bridge/resource"
  "tractor.dev/apptron/bridge/platform/linux"
)

type Menu struct {
  menu
  Menu linux.Menu
}

func New(items []Item) *Menu {
  menu := &Menu{
    menu: menu{
      Handle: resource.NewHandle(),
      Items:  items,
    },
  }

  menu.Menu = createMenu(items)

  return menu
}

func (m *Menu) Destroy() {
  m.Menu.Destroy()
}

func (m *Menu) Popup() int {
  return 0
}

func createMenu(items []Item) linux.Menu {
  menu := linux.Menu_New()

  if menu.Handle != nil {
    for _, it := range items {
      // @Incomplete: accelerators
      item := linux.MenuItem_New(it.ID, it.Title, it.Disabled, it.Selected, it.Separator)

      if !it.Disabled && len(it.SubMenu) > 0 {
        submenu := createMenu(it.SubMenu)
        item.SetSubmenu(submenu)
      }

      menu.AppendItem(item)
    }
  }

  return menu
}
