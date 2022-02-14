package app

/*
#include "../../lib/hostbridge.h"
*/
import "C"

import (
	"os"
	"unsafe"
)

import (
	"github.com/progrium/hostbridge/bridge/menu"
	"github.com/progrium/hostbridge/bridge/window"
)

type Callback func(event Event)

type EventType int

const (
	EventNone       EventType = iota
	EventClose
	EventDestroyed
	EventFocused
	EventResized
	EventMoved
	EventMenuItem
)

func (e EventType) String() string {
	return []string{"none", "close", "destroyed", "focused", "resized", "moved", "menu-item"}[e]
}

type Event struct {
	Type       EventType
	Name       string
	WindowID   window.Handle
	Position   window.Position
	Size       window.Size
	MenuID     uint16
}

type Module struct {
	shouldQuit  bool
	menu        menu.Menu
}

var module Module

func init() {
	module.shouldQuit = false
}

var userMainLoop Callback

//export go_app_main_loop
func go_app_main_loop(data C.Event) {
	if (module.shouldQuit) {
		return
	}

	if (userMainLoop != nil) {
		event := Event{}
		event.Type     = EventType(data.event_type)
		event.Name     = event.Type.String()
		event.WindowID = window.Handle(data.window_id)
		event.Position = window.Position{ X: float64(data.position.x), Y: float64(data.position.y) }
		event.Size     = window.Size{ Width: float64(data.size.width), Height: float64(data.size.height) }
		event.MenuID   = uint16(data.menu_id)

		userMainLoop(event)
	}
}

func Run(callback Callback) {
	if (callback != nil) {
		userMainLoop = callback
		eventLoop := *(*C.EventLoop)(unsafe.Pointer(&window.EventLoop))
		C.run(eventLoop, C.closure(C.go_app_main_loop))
	}
}

func Quit() {
	if (!module.shouldQuit) {
		module.shouldQuit = true

		os.Exit(0)

		// @Incomplete: ideally this would return execution to the main thread
		// but it seems fine because ControlFlow::Exit actually quits the whole process...
		//
		// @MemoryLeak: window.EventLoop destructor needs to be called here if we return execution to go main
	}
}

func Menu() *menu.Menu {
	if (menu.AppMenuWasSet) {
		return &menu.AppMenu
	}

	return nil
}

func SetMenu(m menu.Menu) {
	menu.AppMenu       = m
	menu.AppMenuWasSet = true
}
