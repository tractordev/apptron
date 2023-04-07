package window

import (
	"errors"
	"log"
	"unsafe"

	"github.com/jchv/go-webview2/pkg/edge"

	"tractor.dev/apptron/bridge/api/menu"
	"tractor.dev/apptron/bridge/event"
	. "tractor.dev/apptron/bridge/platform/win32"
	"tractor.dev/apptron/bridge/resource"
)

type Window struct {
	window

	hwnd    HWND
	webview *edge.Chromium

	minSize   POINT
	maxSize   POINT
	placement WINDOWPLACEMENT

	hasMenu       BOOL
	scale         Size
	isTransparent bool

	ID string
}

func init() {
}

func windowCallback(hwnd HWND, message uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	w := (*Window)(unsafe.Pointer(GetWindowLongPtrW(hwnd, GWLP_USERDATA)))

	if w == nil {
		return DefWindowProc(hwnd, message, wParam, lParam)
	}

	switch message {
	case WM_CLOSE:
		event.Emit(event.Event{
			Type:   event.Close,
			Window: w.Handle,
		})

		// @Incomplete: should this still close the window or should that be up to the user?
		// Return 0 to "consume" the close event and prevent the window from closing.
		//return 0

	case WM_DESTROY:
		event.Emit(event.Event{
			Type:   event.Destroyed,
			Window: w.Handle,
		})

	case WM_SETFOCUS:
		event.Emit(event.Event{
			Type:   event.Focused,
			Window: w.Handle,
		})

	case WM_KILLFOCUS:
		event.Emit(event.Event{
			Type:   event.Blurred,
			Window: w.Handle,
		})

	case WM_SIZE:
		if w.webview != nil {
			w.webview.Resize()
		}

		event.Emit(event.Event{
			Type:   event.Resized,
			Window: w.Handle,
			Size:   w.GetInnerSize(),
		})

	case WM_ACTIVATE:
		if w.webview != nil {
			w.webview.Focus()
		}

	case WM_PAINT:
		if w.isTransparent {
			clientRect := RECT{}
			GetClientRect(hwnd, &clientRect)

			ps := PAINTSTRUCT{}
			hdc := BeginPaint(hwnd, &ps)

			bgRgn := CreateRectRgnIndirect(&clientRect)
			hBrush := CreateSolidBrush(0x000000ff)

			FillRgn(hdc, bgRgn, hBrush)

			DeleteObject(HANDLE(bgRgn))
			DeleteObject(HANDLE(hBrush))

			EndPaint(hwnd, &ps)
			return 0
		}

	case WM_MOVE, WM_MOVING:
		if w.webview != nil {
			w.webview.NotifyParentWindowPositionChanged()
		}

		event.Emit(event.Event{
			Type:     event.Moved,
			Window:   w.Handle,
			Position: w.GetOuterPosition(),
		})

	case WM_COMMAND:
		id := LOWORD(uint32(wParam))
		// @Incomplete: can other things trigger WM_COMMAND other than our menu?
		event.Emit(event.Event{
			Type:     event.MenuItem,
			Window:   w.Handle,
			MenuItem: int(id),
		})

	case WM_GETMINMAXINFO:
		info := (*MINMAXINFO)(unsafe.Pointer(lParam))

		style := DWORD(GetWindowLongW(w.hwnd, GWL_STYLE))

		if w.minSize.X > 0 && w.minSize.Y > 0 {
			minSize := mulSize(Size{Width: float64(w.minSize.X), Height: float64(w.minSize.Y)}, w.scale)
			minSize = windowSizeForClientSize(style, minSize, w.hasMenu)
			info.PtMinTrackSize = POINT{X: LONG(minSize.Width), Y: LONG(minSize.Height)}
		}

		if w.maxSize.X > 0 && w.maxSize.Y > 0 {
			maxSize := mulSize(Size{Width: float64(w.maxSize.X), Height: float64(w.maxSize.Y)}, w.scale)
			maxSize = windowSizeForClientSize(style, maxSize, w.hasMenu)
			info.PtMaxTrackSize = POINT{X: LONG(maxSize.Width), Y: LONG(maxSize.Height)}
		}

		return 0

	case WM_DPICHANGED:
		scalex := float64(HIWORD(uint32(wParam))) / float64(USER_DEFAULT_SCREEN_DPI)
		scaley := float64(LOWORD(uint32(wParam))) / float64(USER_DEFAULT_SCREEN_DPI)
		w.scale = Size{Width: scalex, Height: scaley}

		// NOTE(nick): adjust the window rect when the DPI scale changes
		// For example, if you go into the "Make everything bigger" section and change the global pixel scale
		suggested := (*RECT)(unsafe.Pointer(lParam))
		SetWindowPos(w.hwnd, HWND_TOP,
			int(suggested.Left),
			int(suggested.Top),
			int(suggested.Right-suggested.Left),
			int(suggested.Bottom-suggested.Top),
			SWP_NOACTIVATE|SWP_NOZORDER)

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
	log.Println("Callback from JavaScript!!", msg)
}

func New(options Options) (*Window, error) {
	apptronClassName := "APPTRON_WINDOW_CLASS"

	icon := HICON(0)
	if len(options.Icon) > 0 {
		icon = CreateIconFromBytes(options.Icon)
	}

	if !didInitWindowClass {
		// NOTE(nick): setting the icon here sets it for the whole application
		if !RegisterWindowClass(apptronClassName, GetModuleHandle(), windowCallback, 0, icon) {
			return nil, ErrRegisterWindowClass
		}

		didInitWindowClass = true
	}

	var style DWORD = WS_OVERLAPPEDWINDOW

	if options.Frameless {
		style = WS_POPUP
	}

	x := int32(options.Position.X)
	y := int32(options.Position.Y)
	w := int32(0)
	h := int32(0)

	hwnd := CreateWindowExW(0, apptronClassName, options.Title, style, x, y, w, h, 0, 0, GetModuleHandle(), 0)
	if hwnd == 0 {
		return nil, ErrCreateWindow
	}

	if options.Transparent {
		SetWindowLongW(hwnd, GWL_EXSTYLE, GetWindowLongW(hwnd, GWL_EXSTYLE)|WS_EX_LAYERED)

		//
		// When specifying an explicit RGB color, the COLORREF value has the following hexadecimal form: 0x00bbggrr
		// https://docs.microsoft.com/en-us/windows/win32/gdi/colorref
		//
		SetLayeredWindowAttributes(hwnd, 0x000000ff, 0, LWA_COLORKEY)

		// Empty region for the blur effect, so the window is fully transparent
		region := CreateRectRgn(0, 0, -1, -1)

		bb := DWM_BLURBEHIND{}
		bb.DwFlags = DWM_BB_ENABLE | DWM_BB_BLURREGION
		bb.FEnable = TRUE
		bb.HRgnBlur = region
		bb.FTransitionOnMaximized = FALSE

		DwmEnableBlurBehindWindow(hwnd, &bb)
		DeleteObject(HANDLE(region))
	}

	var hasMenu BOOL = FALSE
	m := menu.Main()
	if m != nil {
		SetMenu(hwnd, m.Menu)
		hasMenu = TRUE
	}

	// NOTE(nick): resize window based on pixel scale
	scale := windowGetPixelScale(hwnd)
	size := mulSize(options.Size, scale)

	// NOTE(nick): set default size
	if size.Width == 0 && size.Height == 0 {
		info := MONITORINFOEX{}
		if GetMonitorInfoW(MonitorFromWindow(hwnd, MONITOR_DEFAULTTOPRIMARY), &info) {
			monitorWidth := info.RcMonitor.Right - info.RcMonitor.Left
			monitorHeight := info.RcMonitor.Bottom - info.RcMonitor.Top

			size.Width = float64(int64(float64(monitorWidth) * 0.8))
			size.Height = float64(int64(float64(monitorHeight) * 0.8))
		}
	}

	// size window
	size = windowSizeForClientSize(style, size, hasMenu)
	SetWindowPos(hwnd, HWND_TOP, 0, 0, int(size.Width), int(size.Height), SWP_NOMOVE|SWP_NOOWNERZORDER)

	chromium := edge.NewChromium()
	//chromium.DataPath = options.DataPath
	chromium.SetPermission(edge.CoreWebView2PermissionKindClipboardRead, edge.CoreWebView2PermissionStateAllow)

	if !chromium.Embed(uintptr(hwnd)) {
		DestroyWindow(hwnd)
		return nil, ErrEmbed
	}

	if options.Transparent {
		controller2 := chromium.GetController().GetICoreWebView2Controller2()
		if controller2 != nil {
			color := edge.COREWEBVIEW2_COLOR{}
			color.R = 255
			color.G = 0
			color.B = 0
			color.A = 0

			controller2.PutDefaultBackgroundColor(color)
		}
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

	win := &Window{
		window: window{
			Handle: resource.NewHandle(),
		},
	}
	resource.Retain(win.Handle, win)

	win.hwnd = hwnd
	win.webview = chromium
	win.hasMenu = hasMenu
	win.scale = scale
	win.minSize = POINT{X: LONG(options.MinSize.Width), Y: LONG(options.MinSize.Height)}
	win.maxSize = POINT{X: LONG(options.MaxSize.Width), Y: LONG(options.MaxSize.Height)}
	win.isTransparent = options.Transparent
	win.ID = options.ID

	chromium.MessageCallback = win.messageCallback
	//chromium.Eval("window.chrome.webview.postMessage('Hello, sir!');")

	SetWindowLongPtrW(hwnd, GWLP_USERDATA, unsafe.Pointer(win))

	if options.Center {
		var rect RECT
		if GetWindowRect(hwnd, &rect) {
			windowWidth := LONG(rect.Right - rect.Left)
			windowHeight := LONG(rect.Bottom - rect.Top)

			info := MONITORINFOEX{}
			if GetMonitorInfoW(MonitorFromWindow(hwnd, MONITOR_DEFAULTTOPRIMARY), &info) {
				monitorWidth := info.RcMonitor.Right - info.RcMonitor.Left
				monitorHeight := info.RcMonitor.Bottom - info.RcMonitor.Top

				cx := int(float64(monitorWidth-windowWidth) * 0.5)
				cy := int(float64(monitorHeight-windowHeight) * 0.5)

				SetWindowPos(hwnd, HWND_TOP, cx, cy, 0, 0, SWP_NOSIZE|SWP_NOOWNERZORDER)
			}
		}
	}

	if options.Fullscreen {
		win.SetFullscreen(true)
	}

	if options.Maximized {
		win.SetMaximized(true)
	}

	/*
		if icon != 0 {
			// NOTE(nick): setting the icon here sets it for the specific window
			//SendMessage(hwnd, WM_SETICON, ICON_SMALL, icon)
			//SendMessage(hwnd, WM_SETICON, ICON_BIG, icon)
		}
	*/

	if options.AlwaysOnTop {
		win.SetAlwaysOnTop(true)
	}

	// Finally, present the window and webview.
	if options.Visible {
		win.SetVisible(true)
	}

	event.Emit(event.Event{
		Type:     event.Created,
		Window:   win.Handle,
		Size:     win.GetInnerSize(),
		Position: win.GetOuterPosition(),
	})

	return win, nil
}

func (w *Window) Destroy() {
	if w.webview != nil {
		w.webview.Release()
		w.webview = nil
	}

	if w.hwnd != 0 {
		DestroyWindow(w.hwnd)
		w.hwnd = 0
	}
}

func (w *Window) Focus() {
	SetFocus(w.hwnd)
}

func (w *Window) SetVisible(visible bool) {
	if visible {
		ShowWindow(w.hwnd, SW_SHOW)
		if w.webview != nil {
			w.webview.Show()
		}
	} else {
		ShowWindow(w.hwnd, SW_HIDE)
		if w.webview != nil {
			w.webview.Hide()
		}
	}
}

func (w *Window) IsVisible() bool {
	// @Robustness: is this the same as NSWindow visible?
	//
	// NOTE(nick): from the Apple docs for NSWindow visible:
	// A Boolean value that indicates whether the window is visible onscreen
	//

	return IsWindowVisible(w.hwnd) && !IsIconic(w.hwnd) && !IsWindowCloaked(w.hwnd)
}

func (w *Window) SetMaximized(maximized bool) {
	if maximized {
		ShowWindow(w.hwnd, SW_MAXIMIZE)
	} else {
		ShowWindow(w.hwnd, SW_NORMAL)
	}
}

func (w *Window) SetMinimized(minimized bool) {
	if minimized {
		ShowWindow(w.hwnd, SW_SHOWMINNOACTIVE)
	} else {
		ShowWindow(w.hwnd, SW_NORMAL)
	}
}

func (w *Window) SetFullscreen(fullscreen bool) {
	hwnd := w.hwnd

	style := GetWindowLongW(hwnd, GWL_STYLE)

	if fullscreen {
		info := MONITORINFOEX{}
		if GetWindowPlacement(hwnd, &w.placement) &&
			GetMonitorInfoW(MonitorFromWindow(hwnd, MONITOR_DEFAULTTOPRIMARY), &info) {

			SetWindowLongW(hwnd, GWL_STYLE, style&(^WS_OVERLAPPEDWINDOW))

			SetWindowPos(
				hwnd,
				HWND_TOP,
				int(info.RcMonitor.Left),
				int(info.RcMonitor.Top),
				int(info.RcMonitor.Right-info.RcMonitor.Left),
				int(info.RcMonitor.Bottom-info.RcMonitor.Top),
				SWP_NOOWNERZORDER|SWP_FRAMECHANGED,
			)
		}

	} else {
		SetWindowLongW(hwnd, GWL_STYLE, style|WS_OVERLAPPEDWINDOW)
		SetWindowPlacement(hwnd, &w.placement)

		// @Robustness: is this necessary?
		var flags UINT = SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER | SWP_NOOWNERZORDER | SWP_FRAMECHANGED
		SetWindowPos(hwnd, 0, 0, 0, 0, 0, flags)
	}
}

// NOTE(nick): "inner" size
func (w *Window) SetSize(size Size) {
	size = mulSize(size, w.scale)

	style := DWORD(GetWindowLongW(w.hwnd, GWL_STYLE))
	size = windowSizeForClientSize(style, size, w.hasMenu)

	SetWindowPos(w.hwnd, HWND_TOP, 0, 0, int(size.Width), int(size.Height), SWP_NOMOVE|SWP_NOOWNERZORDER)
}

func (w *Window) SetMinSize(size Size) {
	w.minSize = POINT{X: LONG(size.Width), Y: LONG(size.Height)}
	// NOTE(nick): re-set window size to let WM_GETMINMAXINFO clamp the window if needed
	windowResetSize(w.hwnd)
}

func (w *Window) SetMaxSize(size Size) {
	w.maxSize = POINT{X: LONG(size.Width), Y: LONG(size.Height)}
	// NOTE(nick): re-set window size to let WM_GETMINMAXINFO clamp the window if needed
	windowResetSize(w.hwnd)
}

func (w *Window) SetResizable(resizable bool) {
	style := GetWindowLongW(w.hwnd, GWL_STYLE)

	// NOTE(nick): WS_THICKFRAME controls the windows resizability
	if resizable {
		SetWindowLongW(w.hwnd, GWL_STYLE, style|WS_THICKFRAME)
	} else {
		SetWindowLongW(w.hwnd, GWL_STYLE, style&(^WS_THICKFRAME))
	}
}

func (w *Window) SetAlwaysOnTop(always bool) {
	if always {
		SetWindowPos(w.hwnd, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE|SWP_NOSIZE)
	} else {
		SetWindowPos(w.hwnd, HWND_NOTOPMOST, 0, 0, 0, 0, SWP_NOMOVE|SWP_NOSIZE)
	}
}

func (w *Window) SetPosition(position Position) {
	SetWindowPos(w.hwnd, HWND_TOP, int(position.X), int(position.Y), 0, 0, SWP_NOSIZE)
}

func (w *Window) SetTitle(title string) {
	SetWindowTextW(w.hwnd, title)
}

func (w *Window) GetOuterPosition() Position {
	result := Position{X: 0, Y: 0}

	var rect RECT
	if GetWindowRect(w.hwnd, &rect) {
		result.X = float64(rect.Left)
		result.Y = float64(rect.Top)
	}

	return result
}

func (w *Window) GetOuterSize() Size {
	result := Size{Width: 0, Height: 0}

	scale := windowGetPixelScale(w.hwnd)

	var rect RECT
	if GetWindowRect(w.hwnd, &rect) {
		result.Width = float64(rect.Right-rect.Left) / float64(scale.Width)
		result.Height = float64(rect.Bottom-rect.Top) / float64(scale.Height)
	}

	return result
}

func (w *Window) GetInnerSize() Size {
	result := Size{Width: 0, Height: 0}

	scale := windowGetPixelScale(w.hwnd)

	var rect RECT
	if GetWindowRect(w.hwnd, &rect) {
		result.Width = float64(rect.Right - rect.Left)
		result.Height = float64(rect.Bottom - rect.Top)

		// NOTE(nick): unadjust client rect
		style := DWORD(GetWindowLongW(w.hwnd, GWL_STYLE))
		wr := RECT{0, 0, 0, 0}
		AdjustWindowRect(&wr, style, w.hasMenu)

		width := wr.Right - wr.Left
		height := wr.Bottom - wr.Top

		result.Width -= float64(width)
		result.Height -= float64(height)

		result.Width /= float64(scale.Width)
		result.Height /= float64(scale.Height)
	}

	return result
}

//
// Helpers
//

func windowSizeForClientSize(style DWORD, size Size, menu BOOL) Size {
	wr := RECT{0, 0, LONG(size.Width), LONG(size.Height)}
	AdjustWindowRect(&wr, style, menu)

	return Size{Width: float64(wr.Right - wr.Left), Height: float64(wr.Bottom - wr.Top)}
}

func windowResetSize(hwnd HWND) {
	var rect RECT
	if GetWindowRect(hwnd, &rect) {
		windowWidth := LONG(rect.Right - rect.Left)
		windowHeight := LONG(rect.Bottom - rect.Top)

		var flags UINT = SWP_NOMOVE | SWP_NOZORDER | SWP_NOOWNERZORDER | SWP_FRAMECHANGED
		SetWindowPos(hwnd, 0, 0, 0, int(windowWidth), int(windowHeight), flags)
	}
}

func windowGetPixelScale(hwnd HWND) Size {
	dc := GetDC(hwnd)
	scalex := GetDeviceCaps(dc, LOGPIXELSX)
	scaley := GetDeviceCaps(dc, LOGPIXELSY)
	ReleaseDC(hwnd, dc)

	return Size{
		Width:  float64(float64(scalex) / float64(USER_DEFAULT_SCREEN_DPI)),
		Height: float64(float64(scaley) / float64(USER_DEFAULT_SCREEN_DPI)),
	}
}

func mulSize(a Size, b Size) Size {
	return Size{Width: a.Width * b.Width, Height: a.Height * b.Height}
}
