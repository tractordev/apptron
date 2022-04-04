package window

import (
	"sync"

	"github.com/progrium/macdriver/cocoa"
	mac "github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/progrium/macdriver/webkit"
	"tractor.dev/hostbridge/bridge/event"
	"tractor.dev/hostbridge/bridge/resource"
)

type Window struct {
	window
	cocoa.NSWindow `json:"-"`
}

var ptrLookup sync.Map

func findWindow(win objc.Object) *Window {
	v, ok := ptrLookup.Load(win.Pointer())
	if ok {
		return v.(*Window)
	}
	var found *Window
	resource.Range(func(v interface{}) bool {
		w, ok := v.(*Window)
		if !ok {
			return true
		}
		if w.Pointer() == win.Pointer() {
			found = w
			ptrLookup.Store(win.Pointer(), w)
			return false
		}
		return true
	})
	return found
}

func init() {
	DelegateClass := objc.NewClass("WindowDelegate", "NSObject")
	DelegateClass.AddMethod("windowDidMove:", func(self, notif objc.Object) {
		if win := findWindow(notif.Get("object")); win != nil {
			event.Emit(event.Event{
				Type:     event.Moved,
				Window:   win.Handle,
				Position: win.GetOuterPosition(),
			})
		}
	})
	DelegateClass.AddMethod("windowDidResize:", func(self, notif objc.Object) {
		if win := findWindow(notif.Get("object")); win != nil {
			event.Emit(event.Event{
				Type:   event.Resized,
				Window: win.Handle,
				Size:   win.GetOuterSize(),
			})
		}
	})
	DelegateClass.AddMethod("windowDidBecomeKey:", func(self, notif objc.Object) {
		if win := findWindow(notif.Get("object")); win != nil {
			event.Emit(event.Event{
				Type:   event.Focused,
				Window: win.Handle,
			})
		}
	})
	DelegateClass.AddMethod("windowDidResignKey:", func(self, notif objc.Object) {
		if win := findWindow(notif.Get("object")); win != nil {
			event.Emit(event.Event{
				Type:   event.Blurred,
				Window: win.Handle,
			})
		}
	})
	DelegateClass.AddMethod("windowShouldClose:", func(sender objc.Object) bool {
		// not sure this is right
		if win := findWindow(sender); win != nil {
			event.Emit(event.Event{
				Type:   event.Close,
				Window: win.Handle,
			})
		}
		return true
	})
	DelegateClass.AddMethod("windowWillClose:", func(self, notif objc.Object) {
		// maybe this isn't what should trigger this event
		if win := findWindow(notif.Get("object")); win != nil {
			event.Emit(event.Event{
				Type:   event.Destroyed,
				Window: win.Handle,
			})
		}
	})
	objc.RegisterClass(DelegateClass)
}

func New(options Options) (*Window, error) {
	frame := mac.Rect(options.Position.X, options.Position.Y, options.Size.Width, options.Size.Height)

	nswin := cocoa.NSWindow_Init(
		frame,
		cocoa.NSTitledWindowMask,
		cocoa.NSBackingStoreBuffered,
		false,
	)
	nswin.Retain()
	nswin.MakeKeyAndOrderFront(nil)

	if options.Center {
		screenRect := cocoa.NSScreen_Main().Frame()
		options.Position.X = (screenRect.Size.Width / 2) - (options.Size.Width / 2)
		options.Position.Y = (screenRect.Size.Height / 2) - (options.Size.Height / 2)
		frame = mac.Rect(options.Position.X, options.Position.Y, options.Size.Width, options.Size.Height)
	}

	wkconfig := webkit.WKWebViewConfiguration_New()
	wkconfig.Preferences().SetValueForKey(mac.True, mac.String("developerExtrasEnabled"))

	wv := webkit.WKWebView_Init(mac.Rect(0, 0, 0, 0), wkconfig)
	// NSViewHeightSizable = 16
	// NSViewWidthSizable = 2
	wv.Set("autoresizingMask:", 16|2)
	if options.HTML != "" {
		url := mac.NSURL_URLWithString_(mac.String("http://localhost"))
		wv.LoadHTMLString_baseURL_(mac.String(options.HTML), url)
	}
	if options.URL != "" {
		req := mac.NSURLRequest_Init(mac.URL(options.URL))
		wv.LoadRequest(req)
	}

	mask := uint(cocoa.NSTitledWindowMask |
		cocoa.NSTitledWindowMask |
		cocoa.NSClosableWindowMask |
		cocoa.NSMiniaturizableWindowMask)
	if options.Frameless {
		mask = cocoa.NSBorderlessWindowMask
	}
	if options.Resizable {
		mask = mask | cocoa.NSResizableWindowMask
	}
	nswin.SetStyleMask(mask)

	if options.Title != "" {
		nswin.SetTitle(options.Title)
	} else {
		nswin.SetMovableByWindowBackground(true)
		nswin.SetTitlebarAppearsTransparent(true)
	}

	if options.Transparent {
		nswin.SetBackgroundColor(cocoa.NSColor_Clear())
		nswin.SetOpaque(false)
		wv.SetOpaque(false)
		wv.SetBackgroundColor(cocoa.NSColor_Clear())
		wv.SetValueForKey(mac.False, mac.String("drawsBackground"))
	}

	nswin.SetContentView(wv)

	if options.AlwaysOnTop {
		nswin.SetLevel(cocoa.NSMainMenuWindowLevel)
	}

	nswin.SetFrameDisplay(frame, true)
	nswin.SetDelegate_(objc.Get("WindowDelegate").Alloc().Init())

	win := &Window{
		window: window{
			Handle: resource.NewHandle(),
		},
		NSWindow: nswin,
	}
	resource.Retain(win.Handle, win)

	return win, nil
}

func (w *Window) Destroy() bool {
	w.NSWindow.Close()
	w.NSWindow.Release()
	return true
}

func (w *Window) Focus() {
	if !w.IsMiniaturized() {
		w.MakeKeyAndOrderFront(nil)
		cocoa.NSApp().ActivateIgnoringOtherApps(true)
	}
}

func (w *Window) SetVisible(visible bool) {
	if visible {
		w.MakeKeyAndOrderFront(nil)
	} else {
		w.OrderOut(nil)
	}
}

func (w *Window) IsVisible() bool {
	return w.NSWindow.IsVisible()
}

func (w *Window) SetMaximized(maximized bool) {
	// TODO: if true and is zoomed, return
	// TODO: https://github.com/tauri-apps/tao/blob/dev/src/platform_impl/macos/util/async.rs#L150
}

func (w *Window) SetMinimized(minimized bool) {
	if w.IsMiniaturized() == minimized {
		return
	}
	if minimized {
		w.Miniaturize_(nil)
	} else {
		w.Deminiaturize_(nil)
	}
}

func (w *Window) SetFullscreen(fullscreen bool) {
	// TODO: https://github.com/tauri-apps/tao/blob/dev/src/platform_impl/macos/window.rs#L784
}

func (w *Window) SetSize(size Size) {
	w.SetContentSize_(mac.NSSize{Width: size.Width, Height: size.Height})
}

func (w *Window) SetMinSize(size Size) {
	w.SetContentMinSize_(mac.NSSize{Width: size.Width, Height: size.Height})
}

func (w *Window) SetMaxSize(size Size) {
	w.SetContentMaxSize_(mac.NSSize{Width: size.Width, Height: size.Height})
}

func (w *Window) SetResizable(resizable bool) {
	// TODO: If fullscreen?
	mask := w.StyleMask()
	if resizable {
		mask = mask | cocoa.NSResizableWindowMask
	} else {
		mask = mask & cocoa.NSResizableWindowMask
	}
	w.SetStyleMask(mask)
}

func (w *Window) SetAlwaysOnTop(always bool) {
	if always {
		w.SetLevel(cocoa.NSFloatingWindowLevel)
	} else {
		w.SetLevel(cocoa.NSNormalWindowLevel)
	}
}

func (w *Window) SetPosition(position Position) {
	w.SetFrameTopLeftPoint_(mac.Point(position.X, position.Y))
}

func (w *Window) SetTitle(title string) {
	w.NSWindow.SetTitle(title)
}

func (w *Window) GetOuterPosition() Position {
	frame := w.Frame()
	return Position{
		X: frame.Origin.X,
		Y: frame.Origin.Y,
	}
}

func (w *Window) GetOuterSize() Size {
	frame := w.Frame()
	return Size{
		Width:  frame.Size.Width,
		Height: frame.Size.Height,
	}
}
