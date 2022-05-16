package menu

import (
  "tractor.dev/apptron/bridge/resource"
  "tractor.dev/apptron/bridge/platform/linux"
)

type Menu struct {
  menu
  linux.MenuHandle
}

func New(items []Item) *Menu {
  menu := &Menu{
    menu: menu{
      Handle: resource.NewHandle(),
      Items:  items,
    },
  }

  menu.MenuHandle = createMenu(items)

  return menu
}

func (m *Menu) Destroy() {
}

func (m *Menu) Popup() int {
  return 0
}

func createMenu(items []Item) linux.MenuHandle {
  //linux.OS_Init()

  menu := linux.MenuNew()

  if menu != 0 {
    for _, it := range items {
      // @Incomplete: accelerators
      item := linux.MenuItemNew(it.ID, it.Title, it.Disabled, it.Selected, it.Separator)

      if !it.Disabled && len(it.SubMenu) > 0 {
        submenu := createMenu(it.SubMenu)
        linux.MenuItemSetSubmenu(item, submenu)
      }

      linux.MenuAppendMenuItem(menu, item)
    }
  }

  return menu
}
