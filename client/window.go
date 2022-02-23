package client

import (
	"context"
	"sync"

	"github.com/progrium/qtalk-go/fn"
)

type Handle int

type Position struct {
	X float64
	Y float64
}

type Size struct {
	Width  float64
	Height float64
}

type WindowOptions struct {
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

type WindowModule struct {
	client  *Client
	windows []*Window
	mu      sync.Mutex
}

func (ws *WindowModule) byID(id Handle) *Window {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	for _, w := range ws.windows {
		if w.ID == id {
			return w
		}
	}
	return nil
}

func (s *WindowModule) New(ctx context.Context, opts WindowOptions) (*Window, error) {
<<<<<<< HEAD
	if len(opts.Icon) > 0 {
		opts.IconSel = s.client.ServeData(opts.Icon)
		opts.Icon = nil
	}
=======
>>>>>>> 6bff410 (refactor bridge packages to export to rpc easily, start client with tests)
	var win Window
	_, err := s.client.Call(ctx, "window.New", fn.Args{opts}, &win)
	if err != nil {
		return nil, err
	}
	win.client = s.client
	s.mu.Lock()
	s.windows = append(s.windows, &win)
	s.mu.Unlock()
	return &win, nil
}

type Window struct {
	client *Client

	ID Handle

	OnMoved     func(event Event)
	OnResized   func(event Event)
	OnClose     func(event Event)
	OnFocused   func(event Event)
	OnDestroyed func(event Event)
}

// Destroy
func (w *Window) Destroy(ctx context.Context) (err error) {
	_, err = w.client.Call(ctx, "window.Destroy", fn.Args{w.ID}, nil)
	return
}

// Focus
func (w *Window) Focus(ctx context.Context) (err error) {
	_, err = w.client.Call(ctx, "window.Focus", fn.Args{w.ID}, nil)
	return
}

// GetOuterPosition
func (w *Window) GetOuterPosition(ctx context.Context) (ret Position, err error) {
	_, err = w.client.Call(ctx, "window.GetOuterPosition", fn.Args{w.ID}, &ret)
	return
}

// GetOuterSize
func (w *Window) GetOuterSize(ctx context.Context) (ret Size, err error) {
	_, err = w.client.Call(ctx, "window.GetOuterSize", fn.Args{w.ID}, &ret)
	return
}

// IsDestroyed
func (w *Window) IsDestroyed(ctx context.Context, size Size) (ret bool, err error) {
	_, err = w.client.Call(ctx, "window.IsDestroyed", fn.Args{w.ID, size}, &ret)
	return
}

// IsVisible
func (w *Window) IsVisible(ctx context.Context) (ret bool, err error) {
	_, err = w.client.Call(ctx, "window.IsVisible", fn.Args{w.ID}, &ret)
	return
}

// SetVisible
func (w *Window) SetVisible(ctx context.Context, visible bool) (err error) {
	_, err = w.client.Call(ctx, "window.SetVisible", fn.Args{w.ID, visible}, nil)
	return
}

// SetMaximized
func (w *Window) SetMaximized(ctx context.Context, maximized bool) (err error) {
	_, err = w.client.Call(ctx, "window.SetMaximized", fn.Args{w.ID, maximized}, nil)
	return
}

// SetMinimized
func (w *Window) SetMinimized(ctx context.Context, minimized bool) (err error) {
	_, err = w.client.Call(ctx, "window.SetMinimized", fn.Args{w.ID, minimized}, nil)
	return
}

// SetFullscreen
func (w *Window) SetFullscreen(ctx context.Context, fullscreen bool) (err error) {
	_, err = w.client.Call(ctx, "window.SetFullscreen", fn.Args{w.ID, fullscreen}, nil)
	return
}

// SetMinSize
func (w *Window) SetMinSize(ctx context.Context, size Size) (err error) {
	_, err = w.client.Call(ctx, "window.SetMinSize", fn.Args{w.ID, size}, nil)
	return
}

// SetMaxSize
func (w *Window) SetMaxSize(ctx context.Context, size Size) (err error) {
	_, err = w.client.Call(ctx, "window.SetMaxSize", fn.Args{w.ID, size}, nil)
	return
}

// SetResizable
func (w *Window) SetResizable(ctx context.Context, resizable bool) (err error) {
	_, err = w.client.Call(ctx, "window.SetResizable", fn.Args{w.ID, resizable}, nil)
	return
}

// SetAlwaysOnTop
func (w *Window) SetAlwaysOnTop(ctx context.Context, always bool) (err error) {
	_, err = w.client.Call(ctx, "window.SetAlwaysOnTop", fn.Args{w.ID, always}, nil)
	return
}

// SetSize
func (w *Window) SetSize(ctx context.Context, size Size) (err error) {
	_, err = w.client.Call(ctx, "window.SetSize", fn.Args{w.ID, size}, nil)
	return
}

// SetPosition
func (w *Window) SetPosition(ctx context.Context, pos Position) (err error) {
	_, err = w.client.Call(ctx, "window.SetPosition", fn.Args{w.ID, pos}, nil)
	return
}

// SetTitle
func (w *Window) SetTitle(ctx context.Context, title string) (err error) {
	_, err = w.client.Call(ctx, "window.SetTitle", fn.Args{w.ID, title}, nil)
	return
}
