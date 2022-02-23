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

	focusedWindowID Handle
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
	Icon        []byte
	URL         string
	HTML        string
	Script      string
}

type EventType int

const (
	EventNone EventType = iota
	EventClose
	EventDestroyed
	EventFocused
	EventBlurred
	EventResized
	EventMoved
	EventMenuItem
	EventShortcut
)

func (e EventType) String() string {
	return []string{"none", "close", "destroy", "focus", "blur", "resize", "move", "menu", "shortcut"}[e]
}

type Event struct {
	Type     EventType
	Name     string
	WindowID Handle
	Position Position
	Size     Size
	MenuID   uint16
	Shortcut string
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

func Focused() *Window {
	return Module.Focused()
}

func (m *module) Focused() *Window {
	if m.focusedWindowID >= 0 {
		return m.FindByID(m.focusedWindowID)
	}

	return nil
}

func (m *module) ProcessEvent(event Event) {
	if event.Type == EventFocused {
		m.focusedWindowID = event.WindowID
	}
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
		always_on_top: toCBool(options.AlwaysOnTop),
		frameless:   toCBool(options.Frameless),
		fullscreen: toCBool(options.Fullscreen),
		size: C.Size{ width: C.double(options.Size.Width), height: C.double(options.Size.Height) },
		min_size: C.Size{ width: C.double(options.MinSize.Width), height: C.double(options.MinSize.Height) },
		max_size: C.Size{ width: C.double(options.MaxSize.Width), height: C.double(options.MaxSize.Height) },
		maximized: toCBool(options.Maximized),
		position: C.Position{ x: C.double(options.Position.X), y: C.double(options.Position.Y) },
		resizable: toCBool(options.Resizable),
		title: C.CString(options.Title),
		transparent: toCBool(options.Transparent),
		visible: toCBool(options.Visible),
		center: toCBool(options.Center),
		icon: C.Icon{data: (*C.uchar)(nil), size: C.int(0)},
		url: C.CString(options.URL),
		html: C.CString(options.HTML),
		script: C.CString(options.Script),
	}

	if len(options.Icon) > 0 {
		opts.icon = C.Icon{data: (*C.uchar)(unsafe.Pointer(&options.Icon[0])), size: C.int(len(options.Icon))}
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

func (m *module) Focus(w *Window) {
	w.Focus()
}

func (it *Window) Focus() {
	C.window_set_focused(C.int(it.ID))
}

func (m *module) SetVisible(w *Window, visible bool) {
	w.SetVisible(visible)
}

func (it *Window) SetVisible(visible bool) {
	C.window_set_visible(C.int(it.ID), toCBool(visible))
}

func (m *module) IsVisible(w *Window) bool {
	return w.IsVisible()
}

func (it *Window) IsVisible() bool {
	result := C.window_is_visible(C.int(it.ID))
	return toBool(result)
}

func (m *module) SetMaximized(w *Window, maximized bool) {
	w.SetMaximized(maximized)
}

func (it *Window) SetMaximized(maximized bool) {
	C.window_set_maximized(C.int(it.ID), toCBool(maximized))
}

func (m *module) SetMinimized(w *Window, minimized bool) {
	w.SetMinimized(minimized)
}

func (it *Window) SetMinimized(minimized bool) {
	C.window_set_minimized(C.int(it.ID), toCBool(minimized))
}

func (m *module) SetFullscreen(w *Window, fullscreen bool) {
	w.SetFullscreen(fullscreen)
}

func (it *Window) SetFullscreen(fullscreen bool) {
	C.window_set_fullscreen(C.int(it.ID), toCBool(fullscreen))
}

func (m *module) SetSize(w *Window, size Size) {
	w.SetSize(size)
}

func (it *Window) SetSize(size Size) {
	arg := C.Size{width: C.double(size.Width), height: C.double(size.Height)}
	C.window_set_size(C.int(it.ID), arg)
}

func (m *module) SetMinSize(w *Window, size Size) {
	w.SetMinSize(size)
}

func (it *Window) SetMinSize(size Size) {
	arg := C.Size{width: C.double(size.Width), height: C.double(size.Height)}
	C.window_set_min_size(C.int(it.ID), arg)
}

func (m *module) SetMaxSize(w *Window, size Size) {
	w.SetMaxSize(size)
}

func (it *Window) SetMaxSize(size Size) {
	arg := C.Size{width: C.double(size.Width), height: C.double(size.Height)}
	C.window_set_max_size(C.int(it.ID), arg)
}

func (m *module) SetResizable(w *Window, resizable bool) {
	w.SetResizable(resizable)
}

func (it *Window) SetResizable(resizable bool) {
	C.window_set_resizable(C.int(it.ID), toCBool(resizable))
}

func (m *module) SetAlwaysOnTop(w *Window, always bool) {
	w.SetAlwaysOnTop(always)
}

func (it *Window) SetAlwaysOnTop(always bool) {
	C.window_set_always_on_top(C.int(it.ID), toCBool(always))
}

func (m *module) SetPosition(w *Window, position Position) {
	w.SetPosition(position)
}

func (it *Window) SetPosition(position Position) {
	arg := C.Position{x: C.double(position.X), y: C.double(position.Y)}
	C.window_set_position(C.int(it.ID), arg)
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
