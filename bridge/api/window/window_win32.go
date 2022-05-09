package window

import (
	"tractor.dev/apptron/bridge/resource"
)

type Window struct {
	window
}

func init() {
}

func New(options Options) (*Window, error) {
	win := &Window{
		window: window{
			Handle: resource.NewHandle(),
		},
	}
	resource.Retain(win.Handle, win)

	return win, nil
}

func (w *Window) Destroy() {
}

func (w *Window) Focus() {
}

func (w *Window) SetVisible(visible bool) {
}

func (w *Window) IsVisible() bool {
	return false
}

func (w *Window) SetMaximized(maximized bool) {
	// TODO: if true and is zoomed, return
	// TODO: https://github.com/tauri-apps/tao/blob/dev/src/platform_impl/macos/util/async.rs#L150
}

func (w *Window) SetMinimized(minimized bool) {
}

func (w *Window) SetFullscreen(fullscreen bool) {
	// TODO: https://github.com/tauri-apps/tao/blob/dev/src/platform_impl/macos/window.rs#L784
}

func (w *Window) SetSize(size Size) {
}

func (w *Window) SetMinSize(size Size) {
}

func (w *Window) SetMaxSize(size Size) {
}

func (w *Window) SetResizable(resizable bool) {
}

func (w *Window) SetAlwaysOnTop(always bool) {
}

func (w *Window) SetPosition(position Position) {
}

func (w *Window) SetTitle(title string) {
}

func (w *Window) GetOuterPosition() Position {
	return Position{
		X: 0,
		Y: 0,
	}
}

func (w *Window) GetOuterSize() Size {
	return Size{
		Width:  0,
		Height: 0,
	}
}
