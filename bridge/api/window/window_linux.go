package window

import (
  "log"
  "sync"

  "tractor.dev/apptron/bridge/event"
  "tractor.dev/apptron/bridge/platform/linux"
  "tractor.dev/apptron/bridge/resource"
)

type Window struct {
  window

  Window  linux.Window
  Webview linux.Webview

  callbackId int

  prevPosition linux.Position
  prevSize     linux.Size
}

var ptrLookup sync.Map

func findWindow(win linux.Window) *Window {
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
    if w.Window.Pointer() == win.Pointer() {
      found = w
      ptrLookup.Store(win.Pointer(), w)
      return false
    }
    return true
  })
  return found
}

func init() {
  linux.OS_Init()

  linux.SetGlobalEventCallback(func(it linux.Event) {

    if win := findWindow(it.Window); win != nil {
      if it.Type == linux.Delete {
        event.Emit(event.Event{
          Type:   event.Destroyed,
          Window: win.Handle,
        })
      }

      if it.Type == linux.Destroy {
        event.Emit(event.Event{
          Type:   event.Close,
          Window: win.Handle,
        })
      }

      if it.Type == linux.Configure {
        if it.Position.X != win.prevPosition.X || it.Position.Y != win.prevPosition.Y {
          event.Emit(event.Event{
            Type:     event.Moved,
            Window:   win.Handle,
            Position: win.GetOuterPosition(),
          })

          win.prevPosition = it.Position
        }

        if it.Size.Width != win.prevSize.Width || it.Size.Height != win.prevSize.Height {
          event.Emit(event.Event{
            Type:   event.Resized,
            Window: win.Handle,
            Size:   win.GetOuterSize(),
          })

          win.prevSize = it.Size
        }
      }

      if it.Type == linux.FocusChange {
        if it.FocusIn {
          event.Emit(event.Event{
            Type:   event.Focused,
            Window: win.Handle,
          })
        } else {
          event.Emit(event.Event{
            Type:   event.Blurred,
            Window: win.Handle,
          })
        }
      }
    }

  })
}

func New(options Options) (*Window, error) {
  win := &Window{
    window: window{
      Handle: resource.NewHandle(),
    },
  }
  resource.Retain(win.Handle, win)

  window := linux.Window_New()

  size := options.Size

  // NOTE(nick): set default size
  if size.Width == 0 && size.Height == 0 {
    monitors := linux.Monitors()
    if len(monitors) > 0 {
      m := monitors[0]

      geom := m.Geometry()
      size.Width = float64(geom.Size.Width) * 0.8
      size.Height = float64(geom.Size.Height) * 0.8
    }
  }

  window.SetSize(int(size.Width), int(size.Height))

  if options.MinSize.Width != 0 || options.MinSize.Height != 0 {
    window.SetMinSize(int(options.MinSize.Width), int(options.MinSize.Height))
  }

  if options.MaxSize.Width != 0 || options.MaxSize.Height != 0 {
    window.SetMaxSize(int(options.MaxSize.Width), int(options.MaxSize.Height))
  }

  if options.Center {
    window.Center()
  } else {
    window.SetPosition(int(options.Position.X), int(options.Position.Y))
  }

  if options.Frameless {
    window.SetDecorated(false)
  }

  if options.Fullscreen {
    window.SetFullscreen(true)
  }

  if options.Maximized {
    window.SetMaximized(true)
  }

  window.SetResizable(options.Resizable)

  if options.Title != "" {
    window.SetTitle(options.Title)
  }

  if options.AlwaysOnTop {
    window.SetAlwaysOnTop(true)
  }

  if len(options.Icon) > 0 {
    window.SetIconFromBytes(options.Icon)
  }

  webview := linux.Webview_New()
  webview.SetSettings(linux.DefaultWebviewSettings())

  myCallback := func(result string) {
    log.Println("Callback from JavaScript!!", result)
  }
  callbackId := webview.RegisterCallback("apptron", myCallback)
  webview.Eval("webkit.messageHandlers.apptron.postMessage(JSON.stringify({ hello: 42 }));")

  window.AddWebview(webview)

  if options.Transparent {
    window.SetTransparent(true)
    webview.SetTransparent(true)
  }

  if options.URL != "" {
    webview.Navigate(options.URL)
  }

  if options.HTML != "" {
    webview.SetHtml(options.HTML, "")
  }

  if options.Script != "" {
    webview.AddScript(options.Script)
  }

  if options.Visible {
    window.Show()
  }

  window.BindEventCallback(0)

  win.Window = window
  win.Webview = webview
  win.callbackId = callbackId

  event.Emit(event.Event{
    Type:   event.Created,
    Window: win.Handle,
  })

  return win, nil
}

func (w *Window) Destroy() {
  if w.callbackId != 0 {
    w.Webview.UnregisterCallback(w.callbackId)
    w.callbackId = 0
  }

  w.Webview.Destroy()
  w.Window.Destroy()
}

func (w *Window) Focus() {
  w.Window.Focus()
}

func (w *Window) SetVisible(visible bool) {
  if visible {
    w.Window.Show()
  } else {
    w.Window.Hide()
  }
}

func (w *Window) IsVisible() bool {
  return w.Window.IsVisible()
}

func (w *Window) SetMaximized(maximized bool) {
  w.Window.SetMaximized(maximized)
}

func (w *Window) SetMinimized(minimized bool) {
  w.Window.SetMinimized(minimized)
}

func (w *Window) SetFullscreen(fullscreen bool) {
  w.Window.SetFullscreen(fullscreen)
}

func (w *Window) SetSize(size Size) {
  w.Window.SetSize(int(size.Width), int(size.Height))
}

func (w *Window) SetMinSize(size Size) {
  w.Window.SetMinSize(int(size.Width), int(size.Height))
}

func (w *Window) SetMaxSize(size Size) {
  w.Window.SetMaxSize(int(size.Width), int(size.Height))
}

func (w *Window) SetResizable(resizable bool) {
  w.Window.SetResizable(resizable)
}

func (w *Window) SetAlwaysOnTop(always bool) {
  w.Window.SetAlwaysOnTop(always)
}

func (w *Window) SetPosition(position Position) {
  w.Window.SetPosition(int(position.X), int(position.Y))
}

func (w *Window) SetTitle(title string) {
  w.Window.SetTitle(title)
}

func (w *Window) GetOuterPosition() Position {
  pos := w.Window.GetPosition()
  return Position{
    X: float64(pos.X),
    Y: float64(pos.Y),
  }
}

func (w *Window) GetOuterSize() Size {
  size := w.Window.GetSize()
  return Size{
    Width:  float64(size.Width),
    Height: float64(size.Height),
  }
}
