# hostbridge API


## window module

### `window.getAllWindows(): Handle[]`

An array of all opened browser windows.


* [BrowserWindow.getAllWindows](https://www.electronjs.org/docs/latest/api/browser-window#browserwindowgetallwindows)


### `window.getFocusedWindow(): Handle|null`

The window that is focused in this application, otherwise returns null.

* [BrowserWindow.getFocusedWindow](https://www.electronjs.org/docs/latest/api/browser-window#browserwindowgetfocusedwindow)

### `window.create(options?: Options): Handle`

* [new BrowserWindow](https://www.electronjs.org/docs/latest/api/browser-window#new-browserwindowoptions)
* [wry::application::window::Window::new](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.new)

### `window.Window`
Object representing a window

##### Fields
* `title: string`
* `transparent: bool`
* `size: Size`
* `position: Point`
* `alwaysOnTop: bool`
* `fullscreen: bool`
* `minSize: Size`
* `maxSize: Size`
* `canResize: bool`
* unavailable in wry
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


##### Events
* `closeRequested()` - The window has been requested to close.
  * [wry::application::event::WindowEvent::CloseRequested](https://docs.rs/wry/latest/wry/application/event/enum.WindowEvent.html#variant.CloseRequested)
* `destroyed()` - The window has been closed/destroyed.
  * [wry::application::event::WindowEvent::Destroyed](https://docs.rs/wry/latest/wry/application/event/enum.WindowEvent.html#variant.Destroyed)
* `focused(focus: bool)` - The window gained or lost focus.
  * [wry::application::event::WindowEvent::Focused](https://docs.rs/wry/latest/wry/application/event/enum.WindowEvent.html#variant.Focused)
* `resized(size: Size)` - The size of the window has changed. Contains the client area’s new dimensions.
  * [wry::application::event::WindowEvent::Resized](https://docs.rs/wry/latest/wry/application/event/enum.WindowEvent.html#variant.Resized)
* `moved(position: Point)` - The position of the window has changed. Contains the window’s new position.
  * [wry::application::event::WindowEvent::Moved](https://docs.rs/wry/latest/wry/application/event/enum.WindowEvent.html#variant.Moved)
* not available in wry?
  * 'show'
  * 'hide'
  * 'maximize'
  * 'unmaximize'
  * 'minimize'
  * 'restore'
  * 'enter-fullscreen'
  * 'leave-fullscreen'

##### Methods
* `destroy()` - Release the window object resource.
* `isDestroyed(): bool` - Has the resource been destroyed.
* `focus()` - Bring the window to front and focus on webview control.
  * [wry::application::window::Window.set_focus](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_focus)
  * [wry::webview::WebView.focus](https://docs.rs/wry/latest/wry/webview/struct.WebView.html#method.focus)
* `setVisible(visible: bool)` - Modifies the window’s visibility.
  * If `false`, this will hide the window. If `true`, this will show the window.
  * [wry::application::window::Window.setVisible](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_visible)
* `isVisible(): bool` - Gets the window’s current vibility state.
  * [wry::application:window::Window.is_visible](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.is_visible)
* !isNormal
* `setMaximized(maximized: bool)` - Sets the window to maximized or back.
  * [wry::application:window::Window.set_maximized](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_maximized)
* `setMinimized(minimized: bool)` - Sets the window to minimized or back.
  * [wry::application:window::Window.set_minimized](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_minimized)
* `setFullscreen(fullscreen: bool)` - Sets the window to fullscreen or back.
  * [wry::application:window::Window.set_fullscreen](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_fullscreen)
  * NOTE: uses borderless with no monitor

* `setSize(size: Size)` - Modifies the inner size of the window.
  * [wry::application:window::Window.set_inner_size](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_inner_size)
* `setMinSize(size: Size)` - Sets a minimum dimension size for the window.
  * [wry::application:window::Window.set_min_inner_size](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_min_inner_size)
* `setMaxSize(size: Size)` - Sets a maximum dimension size for the window.
  * [wry::application:window::Window.set_max_inner_size](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_max_inner_size)
* `setResizable(resizable: bool)` - Sets whether the window is resizable or not.
  * [wry::application:window::Window.set_resizable](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_resizable)
* `setAlwaysOnTop(always: bool)` - Change whether or not the window will always be on top of other windows.
  * [wry::application:window::Window.set_always_on_top](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_always_on_top)

* `setPosition(position: Point)` - Modifies the position of the window.
  * [wry::application:window::Window.set_outer_position](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_outer_position)
  * [BrowserWindow#setPosition](https://www.electronjs.org/docs/latest/api/browser-window#winsetpositionx-y-animate)
* `setTitle(title: string)` - Modifies the title of the window.
  * [wry::application:window::Window.set_title](https://docs.rs/wry/latest/wry/application/window/struct.Window.html#method.set_title)
* TODO?:
  * Most of these do not map to wry APIs, but are possible via direct platform calls. Some are convenience wrappers like center
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

### `window.Options`
Options used by `window.create`.

* `alwaysOnTop: bool`
* `frameless: bool` (!with_decorations)
* `fullscreen: bool` (Fullscreen::Borderless(None))
* `size: Size`
* `minSize: Size`
* `maxSize: Size`
* `maximized: bool`
* `position: Point`
* `resizable: bool`
* `title: string`
* `transparent: bool` (both window and webview)
* `visible: bool`
* `center: bool` (convenience)
* `url: string`
  * [wry::webview::WebViewBuilder.with_url](https://docs.rs/wry/latest/wry/webview/struct.WebViewBuilder.html#method.with_url)
* `html: string`
  * [wry::webview::WebViewBuilder.with_html](https://docs.rs/wry/latest/wry/webview/struct.WebViewBuilder.html#method.with_html)
* `script: string`
  * [wry::webview::WebViewBuilder.with_initialization_script](https://docs.rs/wry/latest/wry/webview/struct.WebViewBuilder.html#method.with_initialization_script)
* icon?? todo
* menu?? todo

See also [wry::application::window::WindowBuilder](https://docs.rs/wry/latest/wry/application/window/struct.WindowBuilder.html)


### `window.Point`
Struct for (logical) screen point location

* `x: float64`
* `y: float64` 

* [Electron: Point](https://www.electronjs.org/docs/latest/api/structures/point)
* [wry: LogicalPosition](https://docs.rs/wry/latest/wry/application/dpi/struct.LogicalPosition.html)


### `window.Size`
Struct for (logical) window size

* `width: float64`
* `height: float64` 


* [Electron: Size](https://www.electronjs.org/docs/latest/api/structures/size)
* [wry: LogicalSize](https://docs.rs/wry/latest/wry/application/dpi/struct.LogicalSize.html)

## app module

## menu module

## notification module

## indicator module

## screen module

## shell module

## hotkey module

## dialog module 

## clipboard module


## Notes

fullscreen in wry/tao has two parameterized modes:
* exclusive with a video mode change
* borderless with optional monitor

currently api supports borderless with current monitor (no monitor specified). TODO: add monitor to fullscreen arguments... figure out how exclusive w/ video mode would be added

electron's concept of "frameless" window is wry window
without "decoration". so frameless = !decoration

this api uses logical size in wry land.

electron window minimizable, maximizable, closable, focusable are not available in wry.