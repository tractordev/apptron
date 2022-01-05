# hostbridge API

## window module

### `window.All(): Handle[]`
Returns array of handles to all opened windows.

* [BrowserWindow.getAllWindows](https://www.electronjs.org/docs/latest/api/browser-window#browserwindowgetallwindows)


### `window.Focused(): Handle|null`
Handle to window that is focused in this application, otherwise returns null.

* [BrowserWindow.getFocusedWindow](https://www.electronjs.org/docs/latest/api/browser-window#browserwindowgetfocusedwindow)

### `window.Create(options?: Options): Handle`
Create a new webview window.

* [new BrowserWindow](https://www.electronjs.org/docs/latest/api/browser-window#new-browserwindowoptions)
* wry: [Window::new](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.new)

### `window.Window` [object]
Object representing a window

```golang
type Window struct {
	Title       string
	Transparent bool
	Size        Size
	Position    Position
	AlwaysOnTop bool
	Fullscreen  bool
	MinSize     Size
	MaxSize     Size
	Resizable   bool
}
```

#### Events
* `CloseRequested()` - The window has been requested to close.
  * wry: [WindowEvent::CloseRequested](https://docs.rs/wry/latest/wry/application/event/enum.WindowEvent.html#variant.CloseRequested)
* `Destroyed()` - The window has been closed/destroyed.
  * wry: [WindowEvent::Destroyed](https://docs.rs/wry/latest/wry/application/event/enum.WindowEvent.html#variant.Destroyed)
* `Focused(focus: bool)` - The window gained or lost focus.
  * wry: [WindowEvent::Focused](https://docs.rs/wry/latest/wry/application/event/enum.WindowEvent.html#variant.Focused)
* `Resized(size: Size)` - The size of the window has changed. Contains the client area’s new dimensions.
  * wry: [WindowEvent::Resized](https://docs.rs/wry/latest/wry/application/event/enum.WindowEvent.html#variant.Resized)
* `Moved(position: Position)` - The position of the window has changed. Contains the window’s new position.
  * wry: [WindowEvent::Moved](https://docs.rs/wry/latest/wry/application/event/enum.WindowEvent.html#variant.Moved)


#### Methods
* `Destroy()` - Release the window object resource.

* `IsDestroyed(): bool` - Has the resource been destroyed.

* `Focus()` - Bring the window to front and focus on webview control.
  * wry: [Window.set_focus](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_focus)
  * wry: [WebView.focus](https://docs.rs/wry/latest/wry/webview/struct.WebView.html#method.focus)

* `SetVisible(visible: bool)` - Modifies the window’s visibility.
  * wry: [Window.setVisible](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_visible)

* `IsVisible(): bool` - Gets the window’s current vibility state.
  * wry: [Window.is_visible](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.is_visible)

* `SetMaximized(maximized: bool)` - Sets the window to maximized or back.
  * wry: [Window.set_maximized](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_maximized)

* `SetMinimized(minimized: bool)` - Sets the window to minimized or back.
  * wry: [Window.set_minimized](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_minimized)

* `SetFullscreen(fullscreen: bool)` - Sets the window to fullscreen or back.
  * wry: [Window.set_fullscreen](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_fullscreen)
  * NOTE: uses borderless with no monitor

* `SetSize(size: Size)` - Modifies the inner size of the window.
  * wry: [Window.set_inner_size](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_inner_size)

* `SetMinSize(size: Size)` - Sets a minimum dimension size for the window.
  * wry: [Window.set_min_inner_size](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_min_inner_size)

* `SetMaxSize(size: Size)` - Sets a maximum dimension size for the window.
  * wry: [Window.set_max_inner_size](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_max_inner_size)

* `SetResizable(resizable: bool)` - Sets whether the window is resizable or not.
  * wry: [Window.set_resizable](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_resizable)

* `SetAlwaysOnTop(always: bool)` - Change whether or not the window will always be on top of other windows.
  * wry: [Window.set_always_on_top](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_always_on_top)

* `SetPosition(position: Position)` - Modifies the position of the window.
  * wry: [Window.set_outer_position](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_outer_position)
  * electron: [BrowserWindow#setPosition](https://www.electronjs.org/docs/latest/api/browser-window#winsetpositionx-y-animate)

* `SetTitle(title: string)` - Modifies the title of the window.
  * wry: [Window.set_title](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_title)



#### Not Yet
Most of these do not map to wry APIs, but are possible via direct platform calls. Some can be convenience wrappers like center 

* methods
  * `isNormal`
  * `setHasShadow`
  * `setOpacity`
  * `setIgnoreMouseEvents`
  * `setParentWindow`
  * `getParentWindow`
  * `getChildWindows`
  * `setAspectRatio` (convenience)
  * `moveAbove`
  * `moveTop`
  * `center` (convenience)
  * `blur`
  * `showInactive`
  * `close`
  * `setCanMove`
  * `setCanMinimize`
  * `setCanMaximize`
  * `setCanFullscreen`
  * `setCanClose`
  * `setBounds`
  * `getBounds`
  * `setContentBounds`
  * `getContentBounds`
  * `getNormalBounds`
* events
  * 'show'
  * 'hide'
  * 'maximize'
  * 'unmaximize'
  * 'minimize'
  * 'restore'
  * 'enter-fullscreen'
  * 'leave-fullscreen'
* fields
  * canMove
  * canMinimize
  * canMaximize
  * canClose
  * canFocus
  * canFullscreen
  * isModal
  * backgroundColor
  * hasShadow
  * opacity
  * titleBar
  * roundedCorners
  * frame
  * parent
  * useContentSize

### `window.Options` [struct]
Options used by `window.Create`.

```golang
type Options struct {
	AlwaysOnTop bool
	Frameless   bool // !with_decorations
	Fullscreen  bool // Fullscreen::Borderless(None)
	Size        Size
	MinSize     Size
	MaxSize     Size
	Maximized   bool
	Position    Position
	Resizable   bool
	Title       string
	Transparent bool // both window and webview
	Visible     bool
	Center      bool // convenience
  Icon        string // bytestream callback
	URL         string
	HTML        string
	Script      string
}
```

* [WebViewBuilder.with_url](https://docs.rs/wry/latest/wry/webview/struct.WebViewBuilder.html#method.with_url)
* [WebViewBuilder.with_html](https://docs.rs/wry/latest/wry/webview/struct.WebViewBuilder.html#method.with_html)
* [WebViewBuilder.with_initialization_script](https://docs.rs/wry/latest/wry/webview/struct.WebViewBuilder.html#method.with_initialization_script)
* [wry::application::window::WindowBuilder](https://docs.rs/wry/latest/wry/application/window/struct.WindowBuilder.html)


### `window.Position` [struct]
Logical window position

```golang
type Position struct {
  X float64
  Y float64
}
```

* electron: [Point](https://www.electronjs.org/docs/latest/api/structures/point)
* wry: [LogicalPosition](https://docs.rs/wry/latest/wry/application/dpi/struct.LogicalPosition.html)


### `window.Size` [struct]
Logical window size

```golang
type Size struct {
	Width  float64
	Height float64
}
```

* electron: [Size](https://www.electronjs.org/docs/latest/api/structures/size)
* wry: [LogicalSize](https://docs.rs/wry/latest/wry/application/dpi/struct.LogicalSize.html)


## menu module

### `menu.AppMenu(): Handle|null`
Returns application if set otherwise null.

* electron: [Menu.getApplicationMenu](https://www.electronjs.org/docs/latest/api/menu#menugetapplicationmenu)


### `menu.SetAppMenu(menu: Handle)`
Set application menu. On Windows and Linux, sets the default menu on windows.

* electron: [Menu.setApplicationMenu](https://www.electronjs.org/docs/latest/api/menu#menusetapplicationmenumenu)

### `menu.Create(items: Item[]): Handle`
Create a new webview window.

### `menu.Menu` [object]
Object representing a menu

```golang
type Menu struct {
	Items []Item
}
```

#### Events

* `ItemClicked(item Item)` - A menu item was clicked.
  * wry: [Event::MenuEvent](https://docs.rs/wry/latest/wry/application/event/enum.Event.html#variant.MenuEvent)

### menu.Item [struct]
Menu item

```golang
type Item struct {
  ID          uint16
  Title       string
  SubMenu     []Item
  Enabled     bool
  Selected    bool
  Accelerator string
}
```

* electron: [MenuItem](https://www.electronjs.org/docs/latest/api/menu-item)
* wry: [MenuItemAttributes](https://docs.rs/wry/latest/wry/application/menu/struct.MenuItemAttributes.html)

## notification module

TODO

## indicator module

TODO

## screen module

TODO

## input module

TODO

## shell module

TODO

## Notes

app module is removed for now.

window icon and menu attributes are not included.

fullscreen in wry/tao has two parameterized modes:
* exclusive with a video mode change
* borderless with optional monitor

currently api supports borderless with current monitor (no monitor specified). TODO: add monitor to fullscreen arguments... figure out how exclusive w/ video mode would be added

electron's concept of "frameless" window is wry window
without "decoration". so frameless = !decoration

this api uses logical size in wry land.

electron window minimizable, maximizable, closable, focusable are not available in wry.