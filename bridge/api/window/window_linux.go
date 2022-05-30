package window

import (
  "log"

  "tractor.dev/apptron/bridge/resource"
  "tractor.dev/apptron/bridge/platform/linux"
)

type Window struct {
  window

  Window linux.Window
  Webview linux.Webview
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

  window := linux.Window_New()
  win.Window = window

  myCallback := func(result string) {
    log.Println(result)
  }

  webview := linux.Webview_New()
  linux.Webview_RegisterCallback(webview, myCallback)
  linux.Webview_SetSettings(webview)
  linux.Window_AddWebview(window, webview)


  if options.Center {
    // @Incomplete
  }

  linux.Window_Show(window)

  win.Window  = window
  win.Webview = webview

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
  linux.Window_SetResizable(w.Window, resizable)
}

func (w *Window) SetAlwaysOnTop(always bool) {
}

func (w *Window) SetPosition(position Position) {
}

func (w *Window) SetTitle(title string) {
  linux.Window_SetTitle(w.Window, title)
}

func (w *Window) GetOuterPosition() Position {
  pos := linux.Window_GetPosition(w.Window)
  return Position{
    X: float64(pos.X),
    Y: float64(pos.Y),
  }
}

func (w *Window) GetOuterSize() Size {
  size := linux.Window_GetSize(w.Window)
  return Size{
    Width:  float64(size.Width),
    Height: float64(size.Height),
  }
}
