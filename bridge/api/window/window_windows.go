package window

import (
	"errors"
	"log"
	"unsafe"

	//"tractor.dev/apptron/bridge/event"
	"github.com/jchv/go-webview2/pkg/edge"
	. "tractor.dev/apptron/bridge/platform/win32"
	"tractor.dev/apptron/bridge/resource"
)

type Window struct {
	window

	Window  HWND
	Webview *edge.Chromium

	MinSize   POINT
	MaxSize   POINT
	Placement WINDOWPLACEMENT
}

func init() {
}

func windowCallback(hwnd HWND, message uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	// @Incomplete: emit events
	// @Incomplete: proper window scaling and WM_DPICHANGED handling

	w := (*Window)(unsafe.Pointer(GetWindowLongPtrW(hwnd, GWLP_USERDATA)))

	if w == nil {
		return DefWindowProc(hwnd, message, wParam, lParam)
	}

	switch message {

	case WM_SIZE:
		if w.Webview != nil {
			w.Webview.Resize()
		}

		/*
			case WM_ACTIVATE:
				if w.Webview != nil {
					w.Webview.Focus()
				}
		*/

		/*
			case WM_MOVE, WM_MOVING:
				w.Webview.NotifyParentWindowPositionChanged()
		*/

		/*
			case WM_GETMINMAXINFO:
				info := (*MINMAXINFO)(unsafe.Pointer(lParam))

				// NOTE(nick): we assume "0" max size means the window can be as big as possible
				if w.MaxSize.X == 0 {
					w.MaxSize.X = LONG_MAX
				}
				if w.MaxSize.Y == 0 {
					w.MaxSize.Y = LONG_MAX
				}

				info.PtMinTrackSize = w.MinSize
				info.PtMaxTrackSize = w.MaxSize

				return 0
		*/

	default:
		return DefWindowProc(hwnd, message, wParam, lParam)
	}

	return DefWindowProc(hwnd, message, wParam, lParam)
}

var didInitWindowClass = false

var (
	ErrRegisterWindowClass = errors.New("Failed to register tray window class")
	ErrCreateWindow        = errors.New("Failed to create window")
	ErrEmbed               = errors.New("Failed to embed chromium browser")
)

func (w *Window) messageCallback(msg string) {
	log.Println("Callback!!!", msg)
}

/*
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

	hwnd := CreateWindowExW(0, apptronClassName, options.Title, WS_OVERLAPPEDWINDOW, 0, 0, 320, 240, 0, 0, GetModuleHandle(), 0)
	if hwnd == 0 {
		return nil, ErrCreateWindow
	}

	//SetWindowLongW(hwnd, GWL_USERDATA, 1)

	ShowWindow(hwnd, SW_SHOW)

	chromium := edge.NewChromium()
	chromium.MessageCallback = win.messageCallback
	//chromium.DataPath = options.DataPath
	chromium.SetPermission(edge.CoreWebView2PermissionKindClipboardRead, edge.CoreWebView2PermissionStateAllow)

	if !chromium.Embed(uintptr(hwnd)) {
		return nil, ErrEmbed
	}

	settings, err := chromium.GetSettings()
	if err == nil {
		settings.PutAreDefaultContextMenusEnabled(true)
		settings.PutAreDevToolsEnabled(true)
	}

	chromium.Navigate("data:text/html, <!doctype html><html><body>Hello, Sailor!</body></html>")

	chromium.Resize()

	win.Window = hwnd
	win.Webview = chromium

	return win, nil
}
*/

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

	// @Incomplete: size is in pixels
	// On MacOS, size is in pixels * window scale

	x := int32(options.Position.X)
	y := int32(options.Position.Y)
	w := int32(options.Size.Width)
	h := int32(options.Size.Height)

	x = 0
	y = 0
	w = 640
	h = 480

	hwnd := CreateWindowExW(0, apptronClassName, options.Title, WS_OVERLAPPEDWINDOW, x, y, w, h, 0, 0, GetModuleHandle(), 0)
	if hwnd == 0 {
		return nil, ErrCreateWindow
	}

	chromium := edge.NewChromium()
	chromium.MessageCallback = win.messageCallback
	//chromium.DataPath = options.DataPath
	chromium.SetPermission(edge.CoreWebView2PermissionKindClipboardRead, edge.CoreWebView2PermissionStateAllow)

	if !chromium.Embed(uintptr(hwnd)) {
		return nil, ErrEmbed
	}

	settings, err := chromium.GetSettings()
	if err == nil {
		settings.PutAreDefaultContextMenusEnabled(true)
		settings.PutAreDevToolsEnabled(true)
	}

	if options.URL != "" {
		chromium.Navigate(options.URL)
	}

	if options.HTML != "" {
		chromium.Navigate("data:text/html, " + options.HTML)
	}

	if options.Script != "" {
		chromium.Eval(options.Script)
	}

	chromium.Resize()

	win.Window = hwnd
	win.Webview = chromium
	win.MinSize = POINT{X: LONG(options.MinSize.Width), Y: LONG(options.MinSize.Height)}
	win.MaxSize = POINT{X: LONG(options.MaxSize.Width), Y: LONG(options.MaxSize.Height)}

	SetWindowLongPtrW(hwnd, GWLP_USERDATA, unsafe.Pointer(win))

	// @Incomplete:
	if options.Transparent {
	}

	if options.Fullscreen {
		win.SetFullscreen(true)
	}

	if options.Maximized {
		win.SetMaximized(true)
	}

	if options.Visible {
		win.SetVisible(true)
	}

	if options.AlwaysOnTop {
		win.SetAlwaysOnTop(true)
	}

	return win, nil
}

func (w *Window) Destroy() {
	w.Webview.Release()
	w.Webview = nil

	DestroyWindow(w.Window)
}

func (w *Window) Focus() {
	SetFocus(w.Window)
}

func (w *Window) SetVisible(visible bool) {
	if visible {
		ShowWindow(w.Window, SW_SHOW)
		w.Webview.Show()
	} else {
		ShowWindow(w.Window, SW_HIDE)
		w.Webview.Hide()
	}
}

func (w *Window) IsVisible() bool {
	// @Incomplete: is this the same as NSWindow visible?

	// NOTE(nick): from the Apple docs for NSWindow visible:
	// A Boolean value that indicates whether the window is visible onscreen

	return IsWindowVisible(w.Window) && !IsIconic(w.Window) && !IsWindowCloaked(w.Window)
}

func (w *Window) SetMaximized(maximized bool) {
	if maximized {
		ShowWindow(w.Window, SW_MAXIMIZE)
	} else {
		ShowWindow(w.Window, SW_NORMAL)
	}
}

func (w *Window) SetMinimized(minimized bool) {
	if minimized {
		ShowWindow(w.Window, SW_SHOWMINNOACTIVE)
	} else {
		ShowWindow(w.Window, SW_NORMAL)
	}
}

func (w *Window) SetFullscreen(fullscreen bool) {
	hwnd := w.Window

	style := GetWindowLongW(hwnd, GWL_STYLE)

	if fullscreen {
		monitorInfo := MONITORINFOEX{}
		hmon := MonitorFromWindow(hwnd, MONITOR_DEFAULTTOPRIMARY)

		if GetWindowPlacement(hwnd, &w.Placement) && GetMonitorInfoW(hmon, &monitorInfo) {

			SetWindowLongW(hwnd, GWL_STYLE, style&(^WS_OVERLAPPEDWINDOW))

			SetWindowPos(
				hwnd,
				HWND_TOP,
				int(monitorInfo.RcMonitor.Left),
				int(monitorInfo.RcMonitor.Top),
				int(monitorInfo.RcMonitor.Right-monitorInfo.RcMonitor.Left),
				int(monitorInfo.RcMonitor.Bottom-monitorInfo.RcMonitor.Top),
				SWP_NOOWNERZORDER|SWP_FRAMECHANGED,
			)
		}

	} else {
		SetWindowLongW(hwnd, GWL_STYLE, style|WS_OVERLAPPEDWINDOW)
		SetWindowPlacement(hwnd, &w.Placement)

		// @Robustness: do we need this?
		var flags UINT = SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER | SWP_NOOWNERZORDER | SWP_FRAMECHANGED
		SetWindowPos(hwnd, 0, 0, 0, 0, 0, flags)
	}
}

func (w *Window) SetSize(size Size) {
	wr := RECT{0, 0, LONG(size.Width), LONG(size.Height)}
	style := UINT(GetWindowLongW(w.Window, GWL_STYLE))
	AdjustWindowRect(&wr, style, FALSE)

	SetWindowPos(w.Window, HWND_TOP, int(wr.Right-wr.Left), int(wr.Bottom-wr.Top), 0, 0, SWP_NOMOVE)
}

func (w *Window) SetMinSize(size Size) {
	w.MinSize = POINT{X: LONG(size.Width), Y: LONG(size.Height)}

	// @Incomplete: will this trigger a WM_GETMINMAXINFO event?
	var flags UINT = SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER | SWP_NOOWNERZORDER | SWP_FRAMECHANGED
	SetWindowPos(w.Window, 0, 0, 0, 0, 0, flags)
}

func (w *Window) SetMaxSize(size Size) {
	w.MaxSize = POINT{X: LONG(size.Width), Y: LONG(size.Height)}

	// @Incomplete: will this trigger a WM_GETMINMAXINFO event?
	var flags UINT = SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER | SWP_NOOWNERZORDER | SWP_FRAMECHANGED
	SetWindowPos(w.Window, 0, 0, 0, 0, 0, flags)
}

func (w *Window) SetResizable(resizable bool) {
	style := GetWindowLongW(w.Window, GWL_STYLE)

	// NOTE(nick): WS_THICKFRAME controls the windows resizability
	if resizable {
		SetWindowLongW(w.Window, GWL_STYLE, style|WS_THICKFRAME)
	} else {
		SetWindowLongW(w.Window, GWL_STYLE, style&(^WS_THICKFRAME))
	}
}

func (w *Window) SetAlwaysOnTop(always bool) {
	if always {
		SetWindowPos(w.Window, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE|SWP_NOSIZE)
	} else {
		SetWindowPos(w.Window, HWND_NOTOPMOST, 0, 0, 0, 0, SWP_NOMOVE|SWP_NOSIZE)
	}
}

func (w *Window) SetPosition(position Position) {
	SetWindowPos(w.Window, HWND_TOP, int(position.X), int(position.Y), 0, 0, SWP_NOSIZE)
}

func (w *Window) SetTitle(title string) {
	SetWindowTextW(w.Window, title)
}

func (w *Window) GetOuterPosition() Position {
	result := Position{X: 0, Y: 0}

	var rect RECT
	if GetWindowRect(w.Window, &rect) {
		result.X = float64(rect.Left)
		result.Y = float64(rect.Top)
	}

	return result
}

func (w *Window) GetOuterSize() Size {
	result := Size{Width: 0, Height: 0}

	var rect RECT
	if GetWindowRect(w.Window, &rect) {
		result.Width = float64(rect.Right - rect.Left)
		result.Height = float64(rect.Bottom - rect.Top)
	}

	return result
}
