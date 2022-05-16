package menu

import (
  //"tractor.dev/apptron/bridge/platform/linux"
  "tractor.dev/apptron/bridge/resource"
)

type Menu struct {
  menu
  //*linux.GtkMenuShell
}

func New(items []Item) *Menu {
  menu := &Menu{
    menu: menu{
      Handle: resource.NewHandle(),
      Items:  items,
    },
  }

  //menu.GtkMenuShell = createMenu(items)

  return menu
}

func (m *Menu) Destroy() {
}

func (m *Menu) Popup() int {
  return 0
}

/*
func createMenu(items []Item) *linux.GtkMenuShell {
  menu := linux.MenuNew()

  if menu != nil {
    for _, it := range items {
      // @Incomplete: accelerators
      item := linux.MenuItemNew(it.ID, it.Title, it.Disabled, it.Selected, it.Separator)

      if !it.Disabled && len(it.SubMenu) > 0 {
        submenu := createMenu(it.SubMenu)
        linux.MenuItemAppendSubmenu(submenu, item)
      }

      linux.MenuAppendMenuItem(menu, item)
    }
  }

  return menu
}
*/