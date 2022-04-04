package window

import (
	"context"

	"github.com/progrium/qtalk-go/rpc"
	"tractor.dev/hostbridge/bridge/misc"
	"tractor.dev/hostbridge/bridge/resource"
)

var (
	Module *module
)

func init() {
	Module = &module{}
}

type module struct{}

func Get(handle resource.Handle) (*Window, error) {
	v, err := resource.Lookup(handle)
	if err != nil {
		return nil, err
	}
	w, ok := v.(*Window)
	if !ok {
		return nil, resource.ErrBadHandle
	}
	return w, nil
}

type window struct {
	Handle      resource.Handle
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
	IconSel     string
	Icon        []byte
	URL         string
	HTML        string
	Script      string
}

type Size = misc.Size
type Position = misc.Position

// func Focused() *Window {
// 	Module.mu.Lock()
// 	defer Module.mu.Unlock()

// 	for _, w := range Module.windows {
// 		if w. {
// 			return &Module.windows[index]
// 		}
// 	}

// 	return nil
// }

// func (m *module) All() (result []*Window) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	for _, w := range Gets {
// 		if !w.destroyed {
// 			result = append(result, w)
// 		}
// 	}

// 	return result
// }

func (m *module) New(options Options, call *rpc.Call) (*Window, error) {
	if options.IconSel != "" {
		var err error
		options.Icon, err = misc.FetchData(context.Background(), call, options.IconSel)
		if err != nil {
			return nil, err
		}
	}

	return New(options)
}

func (m *module) Destroy(h resource.Handle) (ret bool, err error) {
	var w *Window
	if w, err = Get(h); err != nil {
		ret = w.Destroy()
		resource.Release(h)
	}
	return
}

func (m *module) IsDestroyed(h resource.Handle) (ret bool, err error) {
	ret = resource.IsReleased(h)
	return
}

func (m *module) Focus(h resource.Handle) (err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		w.Focus()
	}
	return
}

func (m *module) SetVisible(h resource.Handle, visible bool) (err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		w.SetVisible(visible)
	}
	return
}

func (m *module) IsVisible(h resource.Handle) (ret bool, err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		ret = w.IsVisible()
	}
	return
}

func (m *module) SetMaximized(h resource.Handle, maximized bool) (err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		w.SetMaximized(maximized)
	}
	return
}

func (m *module) SetMinimized(h resource.Handle, minimized bool) (err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		w.SetMinimized(minimized)
	}
	return
}

func (m *module) SetFullscreen(h resource.Handle, fullscreen bool) (err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		w.SetFullscreen(fullscreen)
	}
	return
}

func (m *module) SetSize(h resource.Handle, size Size) (err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		w.SetSize(size)
	}
	return
}

func (m *module) SetMinSize(h resource.Handle, size Size) (err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		w.SetMinSize(size)
	}
	return
}

func (m *module) SetMaxSize(h resource.Handle, size Size) (err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		w.SetMaxSize(size)
	}
	return
}

func (m *module) SetResizable(h resource.Handle, resizable bool) (err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		w.SetResizable(resizable)
	}
	return
}

func (m *module) SetAlwaysOnTop(h resource.Handle, always bool) (err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		w.SetAlwaysOnTop(always)
	}
	return
}

func (m *module) SetPosition(h resource.Handle, position Position) (err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		w.SetPosition(position)
	}
	return
}

func (m *module) SetTitle(h resource.Handle, title string) (err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		w.SetTitle(title)
	}
	return
}

func (m *module) GetOuterPosition(h resource.Handle) (ret Position, err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		ret = w.GetOuterPosition()
	}
	return
}

func (m *module) GetOuterSize(h resource.Handle) (ret Size, err error) {
	var w *Window
	if w, err = Get(h); err == nil {
		ret = w.GetOuterSize()
	}
	return
}
