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

	windows []Window
	event_loop C.Event_Loop
	should_quit bool
}

type Window struct {
	Id          int
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

func _EventName(event_type int) string {
	if (event_type == Event_Type__None)      { return "none"; }
	if (event_type == Event_Type__Close)     { return "close"; }
	if (event_type == Event_Type__Destroyed) { return "destroyed"; }
	if (event_type == Event_Type__Focused)   { return "focused"; }
	if (event_type == Event_Type__Resized)   { return "resized"; }
	if (event_type == Event_Type__Moved)     { return "moved"; }
	return "";
}

type Event struct {
	Type       int
	Name       string
	WindowId   int
	Dim        Position
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
	module.event_loop = C.create_event_loop()
	module.should_quit = false
}

func All() (result []Window) {
	/*
	module.mu.Lock()
	defer module.mu.Unlock()
	*/

	result = make([]Window, 0)

	for _, it := range module.windows {
		if (!it.was_destroyed) {
			result = append(result, it)
		}
	}

	return result
}

func FindIndexById(window_id int) int {
	/*
	module.mu.Lock()
	defer module.mu.Unlock()
	*/

	var result int = -1

	for i, v := range module.windows {
    if v.Id == window_id {
    	result = i
    	break
    }
	}

	return result
}

func FindById(window_id int) *Window {
	index := FindIndexById(window_id)
	if (index >= 0) {
		return &module.windows[index]
	}
	return nil
}

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
		result.Name = _EventName(result.Type)
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

		os.Exit(0)

		// @Incomplete: ideally this would return execution to the main thread
		// but it seems fine because ControlFlow::Exit actually quits the whole process...
		//
		// @MemoryLeak: event_loop destructor needs to be called here if we return execution to go main
	}
}

func Create(options Options) (*Window, error) {
	c_options := C.Window_Options{
		transparent: _CBool(options.Transparent),
		decorations: _CBool(!options.Frameless),
    html: C.CString(options.HTML),
	};

	result := C.create_window(module.event_loop, c_options)
	id := int(result)

	window := Window{}
	window.Id = id
	window.Transparent = options.Transparent

	if (id >= 0) {
		module.windows = append(module.windows, window)
		return &window, nil
	}

	return nil, errors.New("Failed to create window")
}

func (it *Window) Destroy() bool {
	result := false

	if (!it.was_destroyed) {
		success := C.destroy_window(C.int(it.Id))
		if (_ToBool(success)) {
			it.was_destroyed = true
			result = true

			index := FindIndexById(it.Id)
			if (index >= 0) {
				module.windows = append(module.windows[:index], module.windows[index+1:]...)
			}
		}
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