package core

/*
#include "../../lib/hostbridge.h"
*/
import "C"

import (
	"os"
	"unsafe"
)

var (
	shouldQuit bool
	eventLoop  C.EventLoop
)

func init() {
	eventLoop = C.create_event_loop()
}

//export go_app_main_loop
func go_app_main_loop(data C.Event) {
	if shouldQuit {
		return
	}

	select {
	case fn := <-dispatchQueue:
		fn()
	default:
	}

	if EventHandler != nil {
		event := Event{}
		event.Type = EventType(data.event_type)
		event.Name = event.Type.String()
		event.WindowID = Handle(data.window_id)
		event.Position = Position{X: float64(data.position.x), Y: float64(data.position.y)}
		event.Size = Size{Width: float64(data.size.width), Height: float64(data.size.height)}
		event.MenuID = uint16(data.menu_id)
		event.Shortcut = C.GoString(data.shortcut)

		EventHandler(event)
	}
}

func EventLoop() unsafe.Pointer {
	return unsafe.Pointer(&eventLoop)
}

func Run(handler func(event Event)) {
	if handler != nil {
		EventHandler = handler
	}
	eventLoop := *(*C.EventLoop)(EventLoop())
	C.run(eventLoop, C.closure(C.go_app_main_loop))
}

func Quit() {
	if !shouldQuit {
		shouldQuit = true

		os.Exit(0)

		// @Incomplete: ideally this would return execution to the main thread
		// but it seems fine because ControlFlow::Exit actually quits the whole process...
		//
		// @MemoryLeak: EventLoop destructor needs to be called here if we return execution to go main
	}
}
