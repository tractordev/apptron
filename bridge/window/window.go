package window

// NOTE: There should be NO space between the comments and the `import "C"` line.

/*
#include "../../lib/hostbridge.h"
*/
import "C"

import (
	"errors"
	"sync"
	"unsafe"

	"github.com/progrium/hostbridge/bridge/menu"
)

type module struct {
	mu sync.Mutex

	windows    []Window
	shouldQuit bool
}

type Position struct {
	X float64
	Y float64
}

type Size struct {
	Width  float64
	Height float64
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

	destroyed bool
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

var EventLoop C.EventLoop
var Module *module

func init() {
	EventLoop = C.create_event_loop()
	Module = &module{}
}

func All() (result []Window) {
	return Module.All()
}

func (m *module) All() (result []Window) {
	/*
		module.mu.Lock()
		defer module.mu.Unlock()
	*/

	for _, it := range m.windows {
		if !it.destroyed {
			result = append(result, it)
		}
	}

	return result
}

func (m *module) FindIndexByID(windowID Handle) int {
	/*
		module.mu.Lock()
		defer module.mu.Unlock()
	*/

	var result int = -1

	for i, v := range m.windows {
		if v.ID == windowID {
			result = i
			break
		}
	}

	return result
}

func (m *module) FindByID(windowID Handle) *Window {
	index := m.FindIndexByID(windowID)
	if index >= 0 {
		return &m.windows[index]
	}
	return nil
}

func Create(options Options) (*Window, error) {
	return Module.Create(options)
}

func (m *module) Create(options Options) (*Window, error) {
	opts := C.Window_Options{
		transparent: toCBool(options.Transparent),
		decorations: toCBool(!options.Frameless),
		html:        C.CString(options.HTML),
	}

	appMenu := *(*C.Menu)(unsafe.Pointer(&menu.AppMenu))
	result := C.window_create(EventLoop, opts, appMenu)
	id := int(result)

	window := Window{}
	window.ID = Handle(id)
	window.Transparent = options.Transparent

	if id >= 0 {
		m.windows = append(m.windows, window)
		return &window, nil
	}

	return nil, errors.New("Failed to create window")
}

func (m *module) Destroy(w *Window) bool {
	return w.Destroy()
}

func (it *Window) Destroy() bool {
	result := false

	if !it.destroyed {
		success := C.window_destroy(C.int(it.ID))
		if toBool(success) {
			it.destroyed = true
			result = true

			index := Module.FindIndexByID(it.ID)
			if index >= 0 {
				Module.windows = append(Module.windows[:index], Module.windows[index+1:]...)
			}
		}
	}

	return result
}

func (m *module) IsDestroyed(w *Window) bool {
	return w.IsDestroyed()
}

func (it *Window) IsDestroyed() bool {
	return it.destroyed
}

func (m *module) SetTitle(w *Window, title string) {
	w.SetTitle(title)
}

func (it *Window) SetTitle(title string) {
	success := C.window_set_title(C.int(it.ID), C.CString(title))
	if toBool(success) {
		it.Title = title
	}
}

func (m *module) SetVisible(w *Window, visible bool) {
	w.SetVisible(visible)
}

func (it *Window) SetVisible(visible bool) {
	C.window_set_visible(C.int(it.ID), toCBool(visible))
}

func (m *module) SetFullscreen(w *Window, fullscreen bool) {
	w.SetFullscreen(fullscreen)
}

func (it *Window) SetFullscreen(fullscreen bool) {
	C.window_set_fullscreen(C.int(it.ID), toCBool(fullscreen))
}

func (m *module) GetOuterPosition(w *Window) Position {
	return w.GetOuterPosition()
}

func (it *Window) GetOuterPosition() Position {
	result := C.window_get_outer_position(C.int(it.ID))
	return Position{X: float64(result.x), Y: float64(result.y)}
}

func (m *module) GetOuterSize(w *Window) Size {
	return w.GetOuterSize()
}

func (it *Window) GetOuterSize() Size {
	result := C.window_get_outer_size(C.int(it.ID))
	return Size{Width: float64(result.width), Height: float64(result.height)}
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
