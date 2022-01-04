# hostbridge API

## window module

#### `window.getAllWindows(): Ref<Window>[]`

An array of all opened browser windows.


* [Electron Docs](https://www.electronjs.org/docs/latest/api/browser-window#browserwindowgetallwindows)


#### `window.getFocusedWindow(): Ref<Window>|null`

The window that is focused in this application, otherwise returns null.

* [Electron Docs](https://www.electronjs.org/docs/latest/api/browser-window#browserwindowgetfocusedwindow)

#### `window.create([options]): Handle<Window>`

* [Electron Docs](https://www.electronjs.org/docs/latest/api/browser-window#new-browserwindowoptions)

#### `window.Window`
Object representing a window

##### Fields
* `id`
* `title`
* url
* icon
* frame
* parent
* isModal
* backgroundColor
* hasShadow
* opacity
* transparent
* titleBar
* roundedCorners
* width
* height
* x
* y
* useContentSize
* center
* alwaysOnTop
* fullscreen
* minWidth
* minHeight
* maxWidth
* maxHeight
* canResize
* canMove
* canMinimize
* canMaximize
* canClose
* canFocus
* canFullscreen

##### Events
* 'close'
* 'closed'
* 'blur'
* 'focus'
* 'show'
* 'hide'
* 'maximize'
* 'unmaximize'
* 'minimize'
* 'restore'
* 'resize'
* 'moved'
* 'enter-fullscreen'
* 'leave-fullscreen'

##### Methods
* `destroy()`
* `close`
* `isDestroyed`
* `focus`
* `blur`
* `show`
* `showInactive`
* `hide`
* `isVisible(): Boolean` - Whether the window is visible to the user.
* `isNormal`
* `maximize`
* `unmaximize`
* `minimize`
* `restore`
* `setURL`
* `setFullScreen`
* `setAspectRatio`
* `setBackgroundColor`
* `setBounds`
* `getBounds`
* `setContentBounds`
* `getContentBounds`
* `getNormalBounds`
* `setSize`
* `setMinimumSize`
* `setMaximumSize`
* `setCanResize`
* `setCanMove`
* `setCanMinimize`
* `setCanMaximize`
* `setCanFullscreen`
* `setCanClose`
* `setAlwaysOnTop`
* `moveAbove`
* `moveTop`
* `center()` - Moves window to the center of the screen.
* `setPosition(x, y)` - Moves window to x and y.
  * `x` Integer
  * `y` Integer
* `setTitle`
* `setHasShadow`
* `setOpacity`
* `setIgnoreMouseEvents`
* `setParentWindow`
* `getParentWindow`
* `getChildWindows`


#### `window.Point`
Struct for screen point location
```go
type Point struct {
  X int
  Y int
}
```

* [Electron Docs](https://www.electronjs.org/docs/latest/api/structures/point)


#### `window.Size`
Struct for window size
```go
type Size struct {
  Width  int
  Height int
}
```

* [Electron Docs](https://www.electronjs.org/docs/latest/api/structures/size)

## app module

## menu module

## notification module

## indicator module

## screen module

## shell module

## hotkey module

## dialog module 

## clipboard module