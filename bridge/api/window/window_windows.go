package window

import (
	"errors"
	"log"

	. "tractor.dev/apptron/bridge/platform/win32"
	"tractor.dev/apptron/bridge/resource"
)

type Window struct {
	window
}

func init() {
}

func windowCallback(hwnd HWND, message uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	// windowID := GetWindowLongW(hwnd, GWL_USERDATA)

	switch message {
	default:
		return DefWindowProc(hwnd, message, wParam, lParam)
	}
}

var didInitWindowClass = false

var (
	ErrRegisterWindowClass = errors.New("Failed to register tray window class!")
)

func New(options Options) (*Window, error) {
	win := &Window{
		window: window{
			Handle: resource.NewHandle(),
		},
	}
	resource.Retain(win.Handle, win)

	apptronClassName := "APPTRON_WINDOW_CLASS"

	if !didInitWindowClass {
		if !RegisterWindowClass(apptronClassName, GetModuleHandle(), windowCallback, CS_HREDRAW|CS_VREDRAW|CS_OWNDC) {
			return nil, ErrRegisterWindowClass
		}

		didInitWindowClass = true
	}

	hwnd, err := CreateWindowExW(0, apptronClassName, "Hello Window", WS_OVERLAPPEDWINDOW, 0, 0, 320, 240, 0, 0, GetModuleHandle(), 0)
	if err != nil {
		log.Println("Failed to create window!", err)
		return nil, err
	}

	SetWindowLongW(hwnd, GWL_USERDATA, 1)

	ShowWindow(hwnd, SW_SHOW)

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
}

func (w *Window) SetMinimized(minimized bool) {
}

func (w *Window) SetFullscreen(fullscreen bool) {
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
