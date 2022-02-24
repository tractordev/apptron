# hostbridge API

## window module

### `window.All(): Window[]`
Returns array of all opened windows.

* [BrowserWindow.getAllWindows](https://www.electronjs.org/docs/latest/api/browser-window#browserwindowgetallwindows)

### `window.Focused(): Window|null`
Returns window that is focused in this application, otherwise returns null.

* [BrowserWindow.getFocusedWindow](https://www.electronjs.org/docs/latest/api/browser-window#browserwindowgetfocusedwindow)

### `window.New(options?: Options): Window`
Create a new webview window.

* [new BrowserWindow](https://www.electronjs.org/docs/latest/api/browser-window#new-browserwindowoptions)
* wry: [Window::new](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.new)


### `window.Options` [struct]
Options used by `window.New`.

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


### `window.Window` [object]
Object representing a window

```golang
type Window struct {

  // options or setters only
	Title       string
	Resizable   bool
	AlwaysOnTop bool
  MinSize     Size
	MaxSize     Size
  Visible     bool // user might be able to change?
  Frameless   bool
	
  // options, setter, or user changeable
  Maximized   bool      // pollable with is_maximized
  Size        Size      // Resized event
	Position    Position  // Moved event
  Fullscreen  bool      // pollable with fullscreen
  Focused     bool      // Focused event
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


## menu module

### `menu.New(items: Item[]): Menu`
Create a new menu.

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
* [accelerators in string format](https://www.electronjs.org/docs/latest/api/accelerator)
* electron: [MenuItem](https://www.electronjs.org/docs/latest/api/menu-item)
* wry: [MenuItemAttributes](https://docs.rs/wry/latest/wry/application/menu/struct.MenuItemAttributes.html)

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


## app module

### `app.Menu(): menu.Menu|null`
Returns application menu if set otherwise null. On Mac this is the menu bar, otherwise this is the default
window menu (which setting is not yet implemented in window module).

* electron: [Menu.getApplicationMenu](https://www.electronjs.org/docs/latest/api/menu#menugetapplicationmenu)
* wry seems to refer to app menu as MenuBar

### `app.SetMenu(menu: menu.Menu)`
Set application menu. On Windows and Linux, sets the default menu on windows.

* electron: [Menu.setApplicationMenu](https://www.electronjs.org/docs/latest/api/menu#menusetapplicationmenumenu)
* wry seems to refer to app menu as MenuBar

### `app.NewIndicator(icon: []byte, menu: menu.Menu)`
Creates a new app indicator / systray icon with menu.

* library TBD?? tauri-runtime? seems to be missing icon data
* much simplified from prev API since limited library support
* wry has a "trayevent" https://docs.rs/wry/0.13.1/wry/application/event/enum.TrayEvent.html

## shell module

### `shell.NewNotification(notification: Notification)`
Creates a desktop notification

 * library: maybe https://docs.rs/notify-rust/4.5.6/notify_rust/
 * much simplified from prev API since limited library support

### `shell.WriteClipboard(text: string)`
Writes text to the clipboard

 * library: https://docs.rs/wry/0.13.1/wry/application/clipboard/index.html

### `shell.ReadClipboard(): string`
Reads text from the clipboard

 * library: https://docs.rs/wry/0.13.1/wry/application/clipboard/index.html

### `shell.ShowFilePicker(picker: FileDialog): string`
Shows the open/save file picker dialog

 * library: https://docs.rs/rfd/latest/rfd/index.html

### `shell.ShowMessage(msg: MessageDialog): bool`
Shows the messagebox dialog, returning a positive button selection

 * library: https://docs.rs/rfd/latest/rfd/index.html

### `shell.RegisterShortcut(accelerator: string): bool`
Registers a global shortcut based on the provided accelorator string and returns
true if successful. Global shortcut presses are signalled by package level event.

 * [accelerators in string format](https://www.electronjs.org/docs/latest/api/accelerator)
 * library: https://docs.rs/wry/0.13.1/wry/application/global_shortcut/index.html

### `shell.IsShortcutRegistered(accelerator: string): bool`
Checks if an accelerator string has been registered as a shortcut.

 * [accelerators in string format](https://www.electronjs.org/docs/latest/api/accelerator)
 * library: https://docs.rs/wry/0.13.1/wry/application/global_shortcut/index.html

### `shell.UnregisterShortcut(accelerator: string)`
Unregisters a shortcut based on accelerator string.

 * [accelerators in string format](https://www.electronjs.org/docs/latest/api/accelerator)
 * library: https://docs.rs/wry/0.13.1/wry/application/global_shortcut/index.html

### `shell.UnregisterShortcut()`
Unregisters all shortcuts registered by this app.

 * library: https://docs.rs/wry/0.13.1/wry/application/global_shortcut/index.html


### Module Events

* `ShortcutPressed(accelerator string)` - A global shortcut was triggered.
  * wry: https://docs.rs/wry/0.13.1/wry/application/event/enum.Event.html#variant.GlobalShortcutEvent


### `shell.Notification` [struct]
A desktop notification

```golang
type Notification struct {
	Title     string
  Subtitle  string
  Body      string
  Silent    bool
}
```

### `shell.FileDialog` [struct]
A file picker/save dialog

```golang
type FileDialog struct {
	Title     string
  Directory string
  Filename  string
  Mode      string    // pickfile, pickfiles, pickfolder, savefile
  Filters   []string  // each string is comma delimited (go,rs,toml)
                      // with optional label prefix (text:go,txt)
}
```

### `shell.MessageDialog` [struct]
A message dialog

```golang
type MessageDialog struct {
	Title   string
  Body    string
  Level   string // info, warning, error
  Buttons string // ok, okcancel, yesno
}
```

## screen module

### `screen.Displays(): Display[]`
Returns array of available displays.

* https://docs.rs/wry/0.13.1/wry/application/monitor/index.html

### `screen.Display` [struct]
A display monitor

```golang
type Display struct {
	Name        string
  Size        Size
  Position    Position
  ScaleFactor float64
}
```
### `screen.Position` [struct]
Logical screen position

```golang
type Position struct {
  X float64
  Y float64
}
```

### `screen.Size` [struct]
Logical screen size

```golang
type Size struct {
	Width  float64
	Height float64
}
```




## Notes

notification and indicator apis simplified for currently available backing libraries

window icon and menu attributes are not included.

fullscreen in wry/tao has two parameterized modes:
* exclusive with a video mode change
* borderless with optional monitor

currently api supports borderless with current monitor (no monitor specified). TODO: add monitor to fullscreen arguments... figure out how exclusive w/ video mode would be added

electron's concept of "frameless" window is wry window
without "decoration". so frameless = !decoration

this api uses logical size in wry land.

electron window minimizable, maximizable, closable, focusable are not available in wry.

#### Window Object, Deferred
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

