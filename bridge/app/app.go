package app

/*
#include "../../lib/hostbridge.h"
*/
import "C"

import (
	"os"
	"unsafe"

	"github.com/progrium/hostbridge/bridge/menu"
	"github.com/progrium/hostbridge/bridge/window"
)

var dispatchQueue chan func()
var module Module
var userMainLoop Callback

func init() {
	module.shouldQuit = false
	dispatchQueue = make(chan func(), 1)
}

type Callback func(event Event)

type EventType int

const (
	EventNone EventType = iota
	EventClose
	EventDestroyed
	EventFocused
	EventResized
	EventMoved
	EventMenuItem
	EventShortcut
)

func (e EventType) String() string {
	return []string{"none", "close", "destroyed", "focused", "resized", "moved", "menu-item", "shortcut"}[e]
}

type Event struct {
	Type     EventType
	Name     string
	WindowID window.Handle
	Position window.Position
	Size     window.Size
	MenuID   uint16
	Shortcut string
}

type Module struct {
	shouldQuit bool
	menu       menu.Menu
}

//export go_app_main_loop
func go_app_main_loop(data C.Event) {
	if module.shouldQuit {
		return
	}

	select {
	case fn := <-dispatchQueue:
		fn()
	default:
	}

	if userMainLoop != nil {
		event := Event{}
		event.Type = EventType(data.event_type)
		event.Name = event.Type.String()
		event.WindowID = window.Handle(data.window_id)
		event.Position = window.Position{X: float64(data.position.x), Y: float64(data.position.y)}
		event.Size = window.Size{Width: float64(data.size.width), Height: float64(data.size.height)}
		event.MenuID = uint16(data.menu_id)
		event.Shortcut = C.GoString(data.shortcut)

		userMainLoop(event)
	}
}

func Dispatch(fn func()) {
	dispatchQueue <- fn
}

func Run(callback Callback) {
	userMainLoop = callback
	eventLoop := *(*C.EventLoop)(unsafe.Pointer(&window.EventLoop))
	C.run(eventLoop, C.closure(C.go_app_main_loop))
}

func Quit() {
	if !module.shouldQuit {
		module.shouldQuit = true

		os.Exit(0)

		// @Incomplete: ideally this would return execution to the main thread
		// but it seems fine because ControlFlow::Exit actually quits the whole process...
		//
		// @MemoryLeak: window.EventLoop destructor needs to be called here if we return execution to go main
	}
}

func Menu() *menu.Menu {
	if menu.AppMenuWasSet {
		return &menu.AppMenu
	}

	return nil
}

func SetMenu(m menu.Menu) {
	menu.AppMenu = m
	menu.AppMenuWasSet = true
}

func NewIndicator(icon []byte, items []menu.Item) {
	eventLoop := *(*C.EventLoop)(unsafe.Pointer(&window.EventLoop))

	var cicon C.Icon
	if len(icon) > 0 {
		cicon = C.Icon{data: (*C.uchar)(unsafe.Pointer(&icon[0])), size: C.int(len(icon))}
	} else {
		cicon = C.Icon{data: (*C.uchar)(nil), size: C.int(0)}
	}

	trayMenu := NewContextMenu(items)

	C.tray_set_system_tray(eventLoop, cicon, trayMenu)
}

func NewContextMenu(items []menu.Item) C.ContextMenu {
	result := C.context_menu_create()

	for _, it := range items {
		if len(it.SubMenu) > 0 {
			submenu := NewContextMenu(it.SubMenu)
			C.context_menu_add_submenu(result, C.CString(it.Title), toCBool(it.Enabled), submenu)
		} else {
			C.context_menu_add_item(result, buildCMenuItem(it))
		}
	}

	return result
}

func buildCMenuItem(item menu.Item) C.Menu_Item {
	return C.Menu_Item{
		id:          C.int(item.ID),
		title:       C.CString(item.Title),
		enabled:     toCBool(item.Enabled),
		selected:    toCBool(item.Selected),
		accelerator: C.CString(item.Accelerator),
	}
}

func toCBool(it bool) C.uchar {
	if it {
		return C.uchar(1)
	}

	return C.uchar(0)
}

func toBool(it C.uchar) bool {
	return int(it) != 0
}
