package window

// NOTE: There should be NO space between the comments and the `import "C"` line.
// The -ldl is necessary to fix the linker errors about `dlsym` that would otherwise appear.

/*
#cgo LDFLAGS: ./lib/libhostbridge.a -ldl -framework Carbon -framework Cocoa -framework CoreFoundation -framework CoreVideo -framework IOKit -framework WebKit
#include "../../lib/hostbridge.h"
*/
import "C"

import (
	"sync"
	"errors"
	"os"
)

type Module struct {
	mu sync.Mutex

	windows     []Window
	eventLoop   C.EventLoop
	shouldQuit  bool
}

type Handle int

type Window struct {
	ID          Handle
	Title       string
	Transparent bool

	/*
	Size        Size
	Position    Position
	AlwaysOnTop bool
	Fullscreen  bool
	MinSize     Size
	MaxSize     Size
	Resizable   bool
	*/

	destroyed   bool
}

type EventType int

const (
    EventNone       EventType = iota
    EventClose
    EventDestroyed
    EventFocused
    EventResized
    EventMoved
)

func (e EventType) String() string {
    return []string{"none", "close", "destroyed", "focused", "resized", "moved"}[e]
}

type Event struct {
	Type       EventType
	Name       string
	WindowID   Handle
	Position   Position
	Size       Size
}

type Options struct {
	Transparent bool
	Frameless   bool
	HTML        string

	/*
	AlwaysOnTop bool
	Size        Size
	MinSize     Size
	MaxSize     Size
	Maximized   bool
	Position    Position
	Resizable   bool
	Title       string
	Transparent bool
	Visible     bool
	Center      bool
	Icon        string // bytestream callback
	URL         string
	Script      string
	*/
}

type Position struct {
	X float64
	Y float64
}

type Size struct {
	Width  float64
	Height float64
}

var module Module

func init() {
	module.eventLoop  = C.create_event_loop()
	module.shouldQuit = false
}

func All() (result []Window) {
	/*
	module.mu.Lock()
	defer module.mu.Unlock()
	*/

	for _, it := range module.windows {
		if (!it.destroyed) {
			result = append(result, it)
		}
	}

	return result
}

func FindIndexByID(windowID Handle) int {
	/*
	module.mu.Lock()
	defer module.mu.Unlock()
	*/

	var result int = -1

	for i, v := range module.windows {
		if v.ID == windowID {
			result = i
			break
		}
	}

	return result
}

func FindByID(windowID Handle) *Window {
	index := FindIndexByID(windowID)
	if (index >= 0) {
		return &module.windows[index]
	}
	return nil
}

type Callback func(event Event)

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
		event.WindowID = Handle(data.window_id)
		event.Position = makePosition(data.position)
		event.Size     = makeSize(data.size)

		userMainLoop(event)
	}
}

func Run(callback Callback) {
	if (callback != nil) {
		userMainLoop = callback
		C.run(module.eventLoop, C.closure(C.go_app_main_loop))
	}
}

func Quit() {
	if (!module.shouldQuit) {
		module.shouldQuit = true

		os.Exit(0)

		// @Incomplete: ideally this would return execution to the main thread
		// but it seems fine because ControlFlow::Exit actually quits the whole process...
		//
		// @MemoryLeak: eventLoop destructor needs to be called here if we return execution to go main
	}
}

func Create(options Options) (*Window, error) {
	opts := C.Window_Options{
		transparent: toCBool(options.Transparent),
		decorations: toCBool(!options.Frameless),
		html: C.CString(options.HTML),
	};

	result := C.window_create(module.eventLoop, opts)
	id := int(result)

	window := Window{}
	window.ID          = Handle(id)
	window.Transparent = options.Transparent

	if (id >= 0) {
		module.windows = append(module.windows, window)
		return &window, nil
	}

	return nil, errors.New("Failed to create window")
}

func (it *Window) Destroy() bool {
	result := false

	if (!it.destroyed) {
		success := C.window_destroy(C.int(it.ID))
		if (toGoBool(success)) {
			it.destroyed = true
			result = true

			index := FindIndexByID(it.ID)
			if (index >= 0) {
				module.windows = append(module.windows[:index], module.windows[index+1:]...)
			}
		}
	}

	return result
}

func (it *Window) IsDestroyed() bool {
	return it.destroyed
}

func (it *Window) SetTitle(title string) {
	success := C.window_set_title(C.int(it.ID), C.CString(title))
	if (toGoBool(success)) {
		it.Title = title
	}
}

func (it *Window) SetVisible(visible bool) {
	C.window_set_visible(C.int(it.ID), toCBool(visible))
}

func (it *Window) SetFullscreen(fullscreen bool) {
	C.window_set_fullscreen(C.int(it.ID), toCBool(fullscreen))
}

func (it *Window) GetOuterPosition() Position {
	result := C.window_get_outer_position(C.int(it.ID))
	return makePosition(result)
}

func (it *Window) GetOuterSize() Size {
	result := C.window_get_outer_size(C.int(it.ID))
	return makeSize(result)
}


func makePosition(it C.Position) Position {
	return Position{ X: float64(it.x), Y: float64(it.y) }
}

func makeSize(it C.Size) Size {
	return Size{ Width: float64(it.width), Height: float64(it.height) }
}

func toCBool(it bool) C.uchar {
	if (it) {
		return C.uchar(1)
	}

	return C.uchar(0)
}

func toGoBool(it C.uchar) bool {
	return int(it) != 0
}