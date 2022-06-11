package window

import (
	"errors"
	"log"

	"github.com/jchv/go-webview2/pkg/edge"
	. "tractor.dev/apptron/bridge/platform/win32"
	"tractor.dev/apptron/bridge/resource"
)

type Window struct {
	window

	Window  HWND
	Webview *edge.Chromium

	MinSize Size
	MaxSize Size
}

func init() {
}

func windowCallback(hwnd HWND, message uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	// windowID := GetWindowLongW(hwnd, GWL_USERDATA)

	switch message {
	case WM_SIZE:
		// chromium.Resize()
	default:
		return DefWindowProc(hwnd, message, wParam, lParam)
	}

	return DefWindowProc(hwnd, message, wParam, lParam)
}

var didInitWindowClass = false

var (
	ErrRegisterWindowClass = errors.New("Failed to register tray window class!")
	ErrEmbed               = errors.New("Failed to embed chromium browser!")
)

func (w *Window) messageCallback(msg string) {
	log.Println("Callback!!!", msg)
}

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
	} else {
		ShowWindow(w.Window, SW_HIDE)
	}
}

func (w *Window) IsVisible() bool {
	// @Incomplete: is this the same as NSWindow visible?
	// Should this also check IsWindowCloaked?

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
	/*
	  HWND hwnd = window->handle;

	  DWORD style = GetWindowLong(hwnd, GWL_STYLE);

	  if (fullscreen) {
	    MONITORINFO monitor_info = {sizeof(monitor_info)};

	    if (
	      GetWindowPlacement(hwnd, &window->placement) &&
	      GetMonitorInfo(MonitorFromWindow(hwnd, MONITOR_DEFAULTTOPRIMARY), &monitor_info)
	    ) {
	      // If the user already had maximized the window, then fullscreening won't work.
	      // So, force the window to be in restored mode before doing the fullscreen
	      //SendMessage(hwnd, WM_SYSCOMMAND, SC_RESTORE, 0);
	      //ShowWindow(hwnd, SW_RESTORE);

	      SetWindowLong(hwnd, GWL_STYLE, style & ~WS_OVERLAPPEDWINDOW);

	      SetWindowPos(
	        hwnd,
	        HWND_TOP,
	        monitor_info.rcMonitor.left,
	        monitor_info.rcMonitor.top,
	        monitor_info.rcMonitor.right  - monitor_info.rcMonitor.left,
	        monitor_info.rcMonitor.bottom - monitor_info.rcMonitor.top,
	        SWP_NOOWNERZORDER | SWP_FRAMECHANGED
	      );
	    }
	  } else {
	    SetWindowLong(hwnd, GWL_STYLE, style | WS_OVERLAPPEDWINDOW);
	    SetWindowPlacement(hwnd, &window->placement);
	    DWORD flags = SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER | SWP_NOOWNERZORDER | SWP_FRAMECHANGED;
	    SetWindowPos(hwnd, 0, 0, 0, 0, 0, flags);
	  }
	*/
}

func (w *Window) SetSize(size Size) {
	SetWindowPos(w.Window, HWND_TOP, int(size.Width), int(size.Height), 0, 0, SWP_NOMOVE)
}

func (w *Window) SetMinSize(size Size) {
	w.MinSize = size

	// @Incomplete: is this enough?
	var flags UINT = SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER | SWP_NOOWNERZORDER | SWP_FRAMECHANGED
	SetWindowPos(w.Window, 0, 0, 0, 0, 0, flags)
}

func (w *Window) SetMaxSize(size Size) {
	w.MaxSize = size

	// @Incomplete: is this enough?
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
