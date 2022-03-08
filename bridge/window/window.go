package window

// NOTE: There should be NO space between the comments and the `import "C"` line.

/*
#include "../../lib/hostbridge.h"
*/
import "C"

import (
	"context"
	"errors"
	"sync"
	"unsafe"

	"github.com/progrium/hostbridge/bridge/app"
	"github.com/progrium/hostbridge/bridge/core"
	"github.com/progrium/qtalk-go/rpc"
)

var (
	Module       *module
	ErrBadHandle = errors.New("bad handle")
)

func init() {
	Module = &module{}
}

type module struct {
	mu sync.Mutex

	windows    []Window
	shouldQuit bool
}

type retVal struct {
	V interface{}
	E error
}

type Window struct {
	ID          core.Handle
	Title       string
	Transparent bool

	/*
		Size        core.Size
		Position    core.Position
		AlwaysOnTop bool
		Fullscreen  bool
		MinSize     core.Size
		MaxSize     core.Size
		Resizable   bool
	*/

	destroyed bool
	mu        sync.Mutex
}

type Options struct {
	AlwaysOnTop bool
	Frameless   bool
	Fullscreen  bool
	Size        core.Size
	MinSize     core.Size
	MaxSize     core.Size
	Maximized   bool
	Position    core.Position
	Resizable   bool
	Title       string
	Transparent bool
	Visible     bool
	Center      bool
	IconSel     string
	Icon        []byte
	URL         string
	HTML        string
	Script      string
}

func All() (result []Window) {
	return Module.All()
}

func (m *module) All() (result []Window) {
	m.mu.Lock()
	defer m.mu.Unlock()

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
	m.mu.Lock()
	defer m.mu.Unlock()

	for index, it := range m.windows {
		is_focused := fromCBool(C.window_is_focused(C.int(it.ID)))
		if is_focused {
			return &m.windows[index]
		}
	}

	return nil
}

func (m *module) FindIndexByID(windowID core.Handle) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result int = -1

	for i, v := range m.windows {
		if v.ID == windowID {
			result = i
			break
		}
	}

	return result
}

func FindByID(windowID core.Handle) *Window {
	return Module.FindByID(windowID)
}

func (m *module) FindByID(windowID core.Handle) *Window {
	index := m.FindIndexByID(windowID)
	m.mu.Lock()
	defer m.mu.Unlock()
	if index >= 0 {
		return &m.windows[index]
	}
	return nil
}

func New(options Options) (*Window, error) {
	opts := C.Window_Options{
		always_on_top: toCBool(options.AlwaysOnTop),
		frameless:     toCBool(options.Frameless),
		fullscreen:    toCBool(options.Fullscreen),
		size:          C.Size{width: C.double(options.Size.Width), height: C.double(options.Size.Height)},
		min_size:      C.Size{width: C.double(options.MinSize.Width), height: C.double(options.MinSize.Height)},
		max_size:      C.Size{width: C.double(options.MaxSize.Width), height: C.double(options.MaxSize.Height)},
		maximized:     toCBool(options.Maximized),
		position:      C.Position{x: C.double(options.Position.X), y: C.double(options.Position.Y)},
		resizable:     toCBool(options.Resizable),
		title:         C.CString(options.Title),
		transparent:   toCBool(options.Transparent),
		visible:       toCBool(options.Visible),
		center:        toCBool(options.Center),
		icon:          C.Icon{data: (*C.uchar)(nil), size: C.int(0)},
		url:           C.CString(options.URL),
		html:          C.CString(options.HTML),
		script:        C.CString(options.Script),
	}
	if len(options.Icon) > 0 {
		opts.icon = C.Icon{data: (*C.uchar)(unsafe.Pointer(&options.Icon[0])), size: C.int(len(options.Icon))}
	}

	appMenu := *(*C.Menu)(unsafe.Pointer(app.Module.Menu()))
	eventLoop := *(*C.EventLoop)(core.EventLoop())
	result := C.window_create(eventLoop, opts, C.int(appMenu.ID))
	id := int(result)

	window := Window{}
	window.ID = core.Handle(id)
	window.Transparent = options.Transparent

	if id >= 0 {
		Module.mu.Lock()
		Module.windows = append(Module.windows, window)
		Module.mu.Unlock()
		return &window, nil
	}
	return nil, errors.New("Failed to create window")
}

func (m *module) New(options Options, call *rpc.Call) (*Window, error) {
	if options.IconSel != "" {
		var err error
		options.Icon, err = core.FetchData(context.Background(), call, options.IconSel)
		if err != nil {
			return nil, err
		}
	}
	ret := make(chan retVal)
	core.Dispatch(func() {
		w, err := New(options)
		ret <- retVal{w, err}
	})
	r := <-ret
	return r.V.(*Window), r.E
}

func (m *module) Destroy(h core.Handle) (bool, error) {
	w := m.FindByID(h)
	if w == nil {
		return false, ErrBadHandle
	}
	ret := make(chan bool)
	core.Dispatch(func() {
		ret <- w.Destroy()
	})
	return <-ret, nil
}

func (w *Window) Destroy() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.destroyed {
		return false
	}

	success := C.window_destroy(C.int(w.ID))
	if !fromCBool(success) {
		return false
	}

	w.destroyed = true

	index := Module.FindIndexByID(w.ID)
	if index >= 0 {
		Module.mu.Lock()
		Module.windows = append(Module.windows[:index], Module.windows[index+1:]...)
		Module.mu.Unlock()
	}
	return true
}

func (m *module) IsDestroyed(h core.Handle) (bool, error) {
	w := m.FindByID(h)
	if w == nil {
		return false, ErrBadHandle
	}
	return w.IsDestroyed(), nil
}

func (w *Window) IsDestroyed() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.destroyed
}

func (m *module) Focus(h core.Handle) error {
	w := m.FindByID(h)
	if w == nil {
		return ErrBadHandle
	}
	core.Dispatch(func() {
		w.Focus()
	})
	return nil
}

func (w *Window) Focus() {
	C.window_set_focused(C.int(w.ID))
}

func (m *module) SetVisible(h core.Handle, visible bool) error {
	w := m.FindByID(h)
	if w == nil {
		return ErrBadHandle
	}
	core.Dispatch(func() {
		w.SetVisible(visible)
	})
	return nil
}

func (it *Window) SetVisible(visible bool) {
	C.window_set_visible(C.int(it.ID), toCBool(visible))
}

func (m *module) IsVisible(h core.Handle) (bool, error) {
	w := m.FindByID(h)
	if w == nil {
		return false, ErrBadHandle
	}
	ret := make(chan bool)
	core.Dispatch(func() {
		ret <- w.IsVisible()
	})
	return <-ret, nil
}

func (w *Window) IsVisible() bool {
	return fromCBool(C.window_is_visible(C.int(w.ID)))
}

func (m *module) SetMaximized(h core.Handle, maximized bool) error {
	w := m.FindByID(h)
	if w == nil {
		return ErrBadHandle
	}
	core.Dispatch(func() {
		w.SetMaximized(maximized)
	})
	return nil
}

func (w *Window) SetMaximized(maximized bool) {
	C.window_set_maximized(C.int(w.ID), toCBool(maximized))
}

func (m *module) SetMinimized(h core.Handle, minimized bool) error {
	w := m.FindByID(h)
	if w == nil {
		return ErrBadHandle
	}
	core.Dispatch(func() {
		w.SetMinimized(minimized)
	})
	return nil
}

func (w *Window) SetMinimized(minimized bool) {
	C.window_set_minimized(C.int(w.ID), toCBool(minimized))
}

func (m *module) SetFullscreen(h core.Handle, fullscreen bool) error {
	w := m.FindByID(h)
	if w == nil {
		return ErrBadHandle
	}
	core.Dispatch(func() {
		w.SetFullscreen(fullscreen)
	})
	return nil
}

func (w *Window) SetFullscreen(fullscreen bool) {
	C.window_set_fullscreen(C.int(w.ID), toCBool(fullscreen))
}

func (m *module) SetSize(h core.Handle, size core.Size) error {
	w := m.FindByID(h)
	if w == nil {
		return ErrBadHandle
	}
	core.Dispatch(func() {
		w.SetSize(size)
	})
	return nil
}

func (w *Window) SetSize(size core.Size) {
	arg := C.Size{width: C.double(size.Width), height: C.double(size.Height)}
	C.window_set_size(C.int(w.ID), arg)
}

func (m *module) SetMinSize(h core.Handle, size core.Size) error {
	w := m.FindByID(h)
	if w == nil {
		return ErrBadHandle
	}
	core.Dispatch(func() {
		w.SetMinSize(size)
	})
	return nil
}

func (w *Window) SetMinSize(size core.Size) {
	arg := C.Size{width: C.double(size.Width), height: C.double(size.Height)}
	C.window_set_min_size(C.int(w.ID), arg)
}

func (m *module) SetMaxSize(h core.Handle, size core.Size) error {
	w := m.FindByID(h)
	if w == nil {
		return ErrBadHandle
	}
	core.Dispatch(func() {
		w.SetMaxSize(size)
	})
	return nil
}

func (w *Window) SetMaxSize(size core.Size) {
	arg := C.Size{width: C.double(size.Width), height: C.double(size.Height)}
	C.window_set_max_size(C.int(w.ID), arg)
}

func (m *module) SetResizable(h core.Handle, resizable bool) error {
	w := m.FindByID(h)
	if w == nil {
		return ErrBadHandle
	}
	core.Dispatch(func() {
		w.SetResizable(resizable)
	})
	return nil
}

func (w *Window) SetResizable(resizable bool) {
	C.window_set_resizable(C.int(w.ID), toCBool(resizable))
}

func (m *module) SetAlwaysOnTop(h core.Handle, always bool) error {
	w := m.FindByID(h)
	if w == nil {
		return ErrBadHandle
	}
	core.Dispatch(func() {
		w.SetAlwaysOnTop(always)
	})
	return nil
}

func (w *Window) SetAlwaysOnTop(always bool) {
	C.window_set_always_on_top(C.int(w.ID), toCBool(always))
}

func (m *module) SetPosition(h core.Handle, position core.Position) error {
	w := m.FindByID(h)
	if w == nil {
		return ErrBadHandle
	}
	core.Dispatch(func() {
		w.SetPosition(position)
	})
	return nil
}

func (w *Window) SetPosition(position core.Position) {
	arg := C.Position{x: C.double(position.X), y: C.double(position.Y)}
	C.window_set_position(C.int(w.ID), arg)
}

func (m *module) SetTitle(h core.Handle, title string) error {
	w := m.FindByID(h)
	if w == nil {
		return ErrBadHandle
	}
	core.Dispatch(func() {
		w.SetTitle(title)
	})
	return nil
}

func (w *Window) SetTitle(title string) {
	if fromCBool(C.window_set_title(C.int(w.ID), C.CString(title))) {
		w.Title = title
	}
}

func (m *module) GetOuterPosition(h core.Handle) (core.Position, error) {
	w := m.FindByID(h)
	if w == nil {
		return core.Position{}, ErrBadHandle
	}
	ret := make(chan core.Position)
	core.Dispatch(func() {
		ret <- w.GetOuterPosition()
	})
	return <-ret, nil
}

func (w *Window) GetOuterPosition() core.Position {
	result := C.window_get_outer_position(C.int(w.ID))
	return core.Position{X: float64(result.x), Y: float64(result.y)}
}

func (m *module) GetOuterSize(h core.Handle) (core.Size, error) {
	w := m.FindByID(h)
	if w == nil {
		return core.Size{}, ErrBadHandle
	}
	ret := make(chan core.Size)
	core.Dispatch(func() {
		ret <- w.GetOuterSize()
	})
	return <-ret, nil
}

func (w *Window) GetOuterSize() core.Size {
	result := C.window_get_outer_size(C.int(w.ID))
	return core.Size{Width: float64(result.width), Height: float64(result.height)}
}

func toCBool(it bool) C.uchar {
	if it {
		return C.uchar(1)
	}

	return C.uchar(0)
}

func fromCBool(it C.uchar) bool {
	return int(it) != 0
}
