package app

import (
  "fmt"
  "os"
  "log"

  "tractor.dev/apptron/bridge/api/menu"
  //"tractor.dev/apptron/bridge/event"
  "tractor.dev/apptron/bridge/platform"
  "tractor.dev/apptron/bridge/platform/linux"
)

var (
  mainMenu *menu.Menu
)

var globalTrayId = 0

func Menu() *menu.Menu {
  return mainMenu
}

func SetMenu(menu *menu.Menu) error {
  //app.SetMainMenu(menu.NSMenu)
  mainMenu = menu
  return nil
}

func NewIndicator(icon []byte, items []menu.Item) {
  //
  // NOTE(nick): it seems like libappindicator warns about the "tmp" directory:
  //
  // libappindicator-WARNING **: 15:49:46.793: Using '/tmp' paths in SNAP environment will lead to unreadable resources
  //
  f, err := os.CreateTemp("", "apptron__icon-*.png")
  if err != nil {
    log.Println("[NewIndicator] Failed to create temporary icon file!")
    return
  }

  _, err = f.Write(icon)
  if err != nil {
    log.Println("[NewIndicator] Failed to create write icon bytes!")
    return
  }

  // @Incomplete @Leak: should remove tmp png file when deleting indicator
  //defer os.Remove(f.Name())

  globalTrayId += 1
  trayId := fmt.Sprintf("tray_%d", globalTrayId)

  menu := menu.New(items)
  linux.NewIndicator(trayId, f.Name(), menu.MenuHandle)
}

func Run(options Options) error {
  platform.Start()
  return nil
}
