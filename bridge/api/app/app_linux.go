package app

import (
  "tractor.dev/apptron/bridge/api/menu"
  //"tractor.dev/apptron/bridge/event"
  "tractor.dev/apptron/bridge/platform"
  "tractor.dev/apptron/bridge/platform/linux"
)

var (
  mainMenu *menu.Menu
)

func init() {
  //linux.OS_Init()
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
  //menu := menu.New(items)
  linux.TestNewIndicator()
}

func Run(options Options) error {
  platform.Start()
  return nil
}
