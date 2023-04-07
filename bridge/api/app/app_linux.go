package app

import (
  "fmt"
  "log"
  "os"
  "sync/atomic"

  "tractor.dev/apptron/bridge/api/menu"
  "tractor.dev/apptron/bridge/event"
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

  trayIconPath := f.Name()

  menu := menu.New(items)
  linux.Indicator_New(trayId, trayIconPath, menu.Menu)

  linux.SetGlobalMenuCallback(func(menuId int) {
    event.Emit(event.Event{
      Type:     event.MenuItem,
      MenuItem: menuId,
    })
  })
}

func Run(options Options) error {
  if options.DisableAutoSave != true {
    setupWindowRestoreListener(options.Identifier)
  }

  // NOTE(nick): MacOS-style window behavior
  if options.Agent == false {
    var windowCount int64

    event.Listen("__APPTRON_Platform_listener__", func(e event.Event) error {
      if e.Type == event.Created {
        atomic.AddInt64(&windowCount, 1)
      }

      if e.Type == event.Destroyed {
        if atomic.AddInt64(&windowCount, -1) == 0 {
          platform.Terminate(true)
        }
      }

      return nil
    })
  }

  platform.Start()
  return nil
}
