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
)

type Module struct {
	handles []int

	mu sync.Mutex

	event_loop C.Event_Loop
	should_quit bool
}

type Window struct {
	Id          int
	Title       string

	/*
	Transparent bool
	Size        Size
	Position    Position
	AlwaysOnTop bool
	Fullscreen  bool
	MinSize     Size
	MaxSize     Size
	Resizable   bool
	*/

	was_destroyed bool
}

const (
  Event_Type__None      int = 0
  Event_Type__Close         = 1
  Event_Type__Destroyed     = 2
  Event_Type__Focused       = 3
  Event_Type__Resized       = 4
  Event_Type__Moved         = 5
)

type Event struct {
	Type       int
	WindowId   int
	Dim        Position
}

type Options struct {
	AlwaysOnTop bool
	Frameless   bool
	Fullscreen  bool
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
	HTML        string
	Script      string
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
	module.event_loop = C.create_event_loop()
	module.should_quit = false
}

/*
func (m *Module) All() (ret []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ret = make([]string, len(m.handles))
	copy(ret, m.handles)
	return ret
}
*/

type User_Callback func(event Event)

var user_main_loop User_Callback

//export go_app_main_loop
func go_app_main_loop(data C.Event) {
	if (module.should_quit) {
		return
	}

	if (user_main_loop != nil) {
		result := Event{}
		result.Type = int(data.event_type)
		result.WindowId = int(data.window_id)
		result.Dim = _MakePosition(data.dim)

		user_main_loop(result)
	}
}

func Run(user_callback User_Callback) {
  if (user_callback != nil) {
	  user_main_loop = user_callback
	  C.run(module.event_loop, C.closure(C.go_app_main_loop))
  }
}

func Quit() {
	if (!module.should_quit) {
		module.should_quit = true

		C.quit(module.event_loop)

		// @MemoryLeak: event_loop destructor needs to be called here
		// but it seems fine because ControlFlow::Exit actually quits the whole process...
	}
}

func Create() (Window, error) {
	result := C.create_window(module.event_loop)
	id := int(result)

	window := Window{}
	window.Id = id

	if (id >= 0) {
		module.handles = append(module.handles, id)
		return window, nil
	}

	return window, errors.New("Failed to create window")
}

func (it *Window) Destroy() bool {
	success := C.destroy_window(C.int(it.Id))
	result := _ToBool(success)

	if (result) {
		it.was_destroyed = true
	}

	return result
}

func (it *Window) IsDestroyed() bool {
	return it.was_destroyed
}

func (it *Window) SetTitle(Title string) {
	success := C.window_set_title(C.int(it.Id), C.CString(Title))
	if (_ToBool(success)) {
		it.Title = Title
	}
}

func (it *Window) SetFullscreen(Fullscreen bool) {
	C.window_set_fullscreen(C.int(it.Id), _CBool(Fullscreen))
}

func (it *Window) GetOuterPosition() Position {
	result := C.window_get_outer_position(C.int(it.Id))
	return _MakePosition(result)
}

func (it *Window) GetOuterSize() Size {
	result := C.window_get_outer_size(C.int(it.Id))
	return _MakeSize(result)
}


func _MakePosition(it C.Vector2) Position {
	return Position{ X: float64(it.x), Y: float64(it.y) }
}

func _MakeSize(it C.Vector2) Size {
	return Size{ Width: float64(it.x), Height: float64(it.y) }
}

func _CBool(it bool) C.uchar {
	if (it) {
		return C.uchar(1)
	}

	return C.uchar(0)
}

func _ToBool(it C.uchar) bool {
	return int(it) != 0
}