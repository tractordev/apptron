package window

// NOTE: There should be NO space between the comments and the `import "C"` line.
// The -ldl is necessary to fix the linker errors about `dlsym` that would otherwise appear.

/*
#cgo LDFLAGS: ./lib/libhostbridge.a -ldl -framework Carbon -framework Cocoa -framework CoreFoundation -framework CoreVideo -framework IOKit -framework WebKit
#include "../../lib/hostbridge.h"
*/
import "C"

import "sync"

type Module struct {
	handles []string

	mu sync.Mutex

	event_loop C.Event_Loop
}

type Window struct {
	Id    int
	Title string
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
}

func (m *Module) All() (ret []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ret = make([]string, len(m.handles))
	copy(ret, m.handles)
	return ret
}

var user_main_loop func(event_type int)

//export go_main_loop
func go_main_loop(i C.int) {
	event_type := int(i)

	if (user_main_loop != nil) {
		user_main_loop(event_type)
	}
}

func Run(user_callback func(event_type int)) {
  if (user_callback != nil) {
	  user_main_loop = user_callback
	  C.run(module.event_loop, C.closure(C.go_main_loop))
  }
}

func Create() (Window, error) {
	result := C.create_window(module.event_loop)
	id := int(result)

	window := Window{}
	window.Id = id

	return window, nil
}

func (it *Window) Destroy() bool {
	result := C.destroy_window(C.int(it.Id))
	return _ToBool(result)
}

func (it *Window) SetTitle(Title string) {
	success := C.window_set_title(C.int(it.Id), C.CString(Title))
	if (_ToBool(success)) {
		it.Title = Title
	}
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

func _ToBool(it C.bool) bool {
	return int(it) == 1
}