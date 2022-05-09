package win32

import (
	"log"
	"syscall"
	"unsafe"
)

//
// Imports
//

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	pGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	pCreateWindowExW     = user32.NewProc("CreateWindowExW")
	pDefWindowProcW      = user32.NewProc("DefWindowProcW")
	pDestroyWindow       = user32.NewProc("DestroyWindow")
	pDispatchMessageW    = user32.NewProc("DispatchMessageW")
	pGetMessageW         = user32.NewProc("GetMessageW")
	pPeekMessageW        = user32.NewProc("PeekMessageW")
	pLoadCursorW         = user32.NewProc("LoadCursorW")
	pPostQuitMessage     = user32.NewProc("PostQuitMessage")
	pRegisterClassExW    = user32.NewProc("RegisterClassExW")
	pTranslateMessage    = user32.NewProc("TranslateMessage")
	pGetCursorPos        = user32.NewProc("GetCursorPos")
	pSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	pGetActiveWindow     = user32.NewProc("GetActiveWindow")

	pCreateMenu       = user32.NewProc("CreateMenu")
	pCreatePopupMenu  = user32.NewProc("CreatePopupMenu")
	pDestroyMenu      = user32.NewProc("DestroyMenu")
	pTrackPopupMenu   = user32.NewProc("TrackPopupMenu")
	pInsertMenuItemW  = user32.NewProc("InsertMenuItemW")
	pGetMenuItemCount = user32.NewProc("GetMenuItemCount")

	pSetProcessDpiAwarenessContext = user32.NewProc("SetProcessDpiAwarenessContext")
)

var (
	shell32 = syscall.NewLazyDLL("shell32.dll")

	pShell_NotifyIconW = shell32.NewProc("Shell_NotifyIconW")
)

func GetModuleHandle() HINSTANCE {
	ret, _, _ := pGetModuleHandleW.Call(uintptr(0))
	return HINSTANCE(ret)
}

func CreateWindow(className, windowName string, style uint32, x, y, width, height int32, parent, menu, instance HINSTANCE) (HWND, error) {
	ret, _, err := pCreateWindowExW.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(className))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(windowName))),
		uintptr(style),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(parent),
		uintptr(menu),
		uintptr(instance),
		uintptr(0),
	)
	if ret == 0 {
		return 0, err
	}
	return HWND(ret), nil
}

func DefWindowProc(hwnd HWND, msg uint32, wparam WPARAM, lparam LPARAM) LRESULT {
	ret, _, _ := pDefWindowProcW.Call(uintptr(hwnd), uintptr(msg), uintptr(wparam), uintptr(lparam))
	return LRESULT(ret)
}

func DestroyWindow(hwnd HWND) error {
	ret, _, err := pDestroyWindow.Call(uintptr(hwnd))
	if ret == 0 {
		return err
	}
	return nil
}

func GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) bool {
	ret, _, _ := pGetMessageW.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
	)
	return int32(ret) != 0
}

func PeekMessageW(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32, removeMsg uint32) bool {
	ret, _, _ := pPeekMessageW.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
		uintptr(removeMsg),
	)
	return int32(ret) != 0
}

func LoadCursorResource(cursorName uint32) (HCURSOR, error) {
	ret, _, err := pLoadCursorW.Call(uintptr(0), uintptr(uint16(cursorName)))
	if ret == 0 {
		return 0, err
	}
	return HCURSOR(ret), nil
}

func TranslateMessage(msg *MSG) {
	pTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
}

func DispatchMessage(msg *MSG) {
	pDispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
}

func PostQuitMessage(exitCode int32) {
	pPostQuitMessage.Call(uintptr(exitCode))
}

func RegisterClassEx(wcx *WNDCLASSEXW) (uint16, error) {
	ret, _, err := pRegisterClassExW.Call(uintptr(unsafe.Pointer(wcx)))
	if ret == 0 {
		return 0, err
	}
	return uint16(ret), nil
}

func SetProcessDpiAwarenessContext(context HANDLE) bool {
	ret, _, _ := pSetProcessDpiAwarenessContext.Call(uintptr(context))
	return ret != 0
}

func SetForegroundWindow(hwnd HWND) bool {
	ret, _, _ := pSetForegroundWindow.Call(uintptr(hwnd))
	return ret != 0
}

func GetActiveWindow() HWND {
	ret, _, _ := pGetActiveWindow.Call()
	return HWND(ret)
}

func GetCursorPos(pos *POINT) bool {
	ret, _, _ := pGetCursorPos.Call(uintptr(unsafe.Pointer(pos)))
	return ret != 0
}

func Shell_NotifyIconW(dwMessage DWORD, nid *NOTIFYICONDATA) bool {
	ret, _, _ := pShell_NotifyIconW.Call(uintptr(dwMessage), uintptr(unsafe.Pointer(nid)))
	return ret != 0
}

func CreatePopupMenu() HMENU {
	ret, _, _ := pCreatePopupMenu.Call()
	return HMENU(ret)
}

func DestroyMenu(menu HMENU) bool {
	ret, _, _ := pDestroyMenu.Call(uintptr(menu))
	return ret != 0
}

func TrackPopupMenu(menu HMENU, flags UINT, x int32, y int32, nReserved int32, hwnd HWND, rect *RECT) int32 {
	result, _, _ := pTrackPopupMenu.Call(
		uintptr(menu),
		uintptr(flags),
		uintptr(x),
		uintptr(y),
		uintptr(nReserved),
		uintptr(hwnd),
		uintptr(unsafe.Pointer(rect)),
	)
	return int32(result)
}

func InsertMenuItemW(menu HMENU, item UINT, byPosition int32, itemInfo *MENUITEMINFO) bool {
	result, _, _ := pInsertMenuItemW.Call(
		uintptr(menu),
		uintptr(item),
		uintptr(byPosition),
		uintptr(unsafe.Pointer(itemInfo)),
	)
	return result != 0
}

func GetMenuItemCount(menu HMENU) int32 {
	result, _, _ := pGetMenuItemCount.Call(uintptr(menu))
	return int32(result)
}

//
// Functions
//

func PollEvents() {
	for {
		msg := MSG{}
    if (!PeekMessageW(&msg, NULL, 0, 0, PM_REMOVE)) {
      break
    }

    /*
    switch (msg.message) {
    case WM_QUIT:
    	PostQuitMessage(0)
    }
    */

		TranslateMessage(&msg)
		DispatchMessage(&msg)
	}
}


func RegisterWindowClass(className string, instance HINSTANCE, callback WNDPROC) bool {
	cursor, err := LoadCursorResource(IDC_ARROW)
	if err != nil {
		log.Println(err)
		return false
	}

	wc := WNDCLASSEXW{
		LpfnWndProc:   syscall.NewCallback(callback),
		HInstance:     instance,
		HCursor:       cursor,
		HbrBackground: COLOR_WINDOW + 1,
		LpszClassName: syscall.StringToUTF16Ptr(className),
	}
	wc.CbSize = UINT(unsafe.Sizeof(wc))

	if _, err = RegisterClassEx(&wc); err != nil {
		log.Println(err)
		return false
	}

	return true
}

func MakeMenuItemSeparator() MENUITEMINFO {
	result := MENUITEMINFO{}
	result.CbSize     = UINT(unsafe.Sizeof(result))
	result.FMask      = MIIM_ID | MIIM_DATA | MIIM_TYPE
	result.WID        = 0
	result.DwItemData = 0
	return result
}

func MakeMenuItem(id int, label string, disabled bool, checked bool, isRadio bool) MENUITEMINFO {
	result := MENUITEMINFO{}

	result.CbSize     = UINT(unsafe.Sizeof(result))
	result.FMask      = MIIM_ID | MIIM_STATE | MIIM_DATA | MIIM_TYPE
	result.FType      = MFT_STRING

	result.FState     = 0
	if checked {
		result.FState |= MFS_CHECKED
	} else {
		result.FState |= MFS_UNCHECKED
	}

	if disabled {
		result.FState |= MFS_DISABLED
	} else {
		result.FState |= MFS_ENABLED
	}

	result.WID        = UINT(id)
	result.DwTypeData = syscall.StringToUTF16Ptr(label)

	if isRadio {
		result.FType |= MFT_RADIOCHECK
	}

	return result
}

func AppendSubmenu(submenu HMENU, mii *MENUITEMINFO) {
	mii.FMask |= MIIM_SUBMENU
	mii.HSubMenu = submenu
}


// NOTE(nick): system tray menu
// @Robustness: add support for multiple tray icons?
var trayIconData NOTIFYICONDATA
var trayWindow   HWND
var trayMenu     HMENU
var trayCallback func(id int32)

const Win32TrayIconMessage = (WM_USER + 1)

func trayWindowCallback(hwnd HWND, message uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	switch message {
		case Win32TrayIconMessage:
			switch lParam {
				case WM_LBUTTONDOWN, WM_RBUTTONDOWN:

					SetForegroundWindow(hwnd)

					mousePosition := POINT{}
					GetCursorPos(&mousePosition)

					result := TrackPopupMenu(trayMenu, TPM_RIGHTBUTTON | TPM_NONOTIFY | TPM_RETURNCMD, int32(mousePosition.X), int32(mousePosition.Y), 0, hwnd, nil)

					if trayCallback != nil {
						trayCallback(result)
					}

				default: break
			}

		default:
			return DefWindowProc(hwnd, message, wParam, lParam)
	}
	return 0
}

func SetTrayMenu(menu HMENU, icon []byte, callback func(id int32)) bool {
	if trayWindow == NULL {
	  trayClassName := "APPTRON_TRAY_WINDOW_CLASS"

	  if (!RegisterWindowClass(trayClassName, GetModuleHandle(), trayWindowCallback)) {
	  	log.Println("Failed to register tray window class!")
	  	return false
	  }

	  hwnd, err := CreateWindow(trayClassName, "Tray Window", 0, 0, 0, 1, 1, 0, 0, GetModuleHandle());
	  if err != nil {
	  	log.Println("Failed to create tray window:", err)
	  	return false
	  }

	  trayWindow = hwnd
	}

	if trayIconData.CbSize > 0 {
		Shell_NotifyIconW(NIM_DELETE, &trayIconData)
	}

	if trayMenu != 0 {
		DestroyMenu(trayMenu)
	}

  trayMenu     = menu
  trayCallback = callback

	trayIconData = NOTIFYICONDATA{}
  trayIconData.CbSize = DWORD(unsafe.Sizeof(trayIconData))
  trayIconData.HWnd = trayWindow
  trayIconData.UID = 0
  trayIconData.UFlags = NIF_MESSAGE | NIF_ICON | NIF_TIP
  trayIconData.UCallbackMessage = Win32TrayIconMessage
  // @Incomplete: figure out icons
  //trayIconData.HIcon = LoadIcon(GetModuleHandle(0), MAKEINTRESOURCE(101));
  trayIconData.SzTip[0] = 0 // @Incomplete: we should put the app name here

  Shell_NotifyIconW(NIM_ADD, &trayIconData)
  return true
}

func testWindowCallback(hwnd HWND, message uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	switch message {
		case WM_CLOSE:
			DestroyWindow(hwnd)
		case WM_DESTROY:
			PostQuitMessage(0)
		default:
			return DefWindowProc(hwnd, message, wParam, lParam)
	}
	return 0
}

func CreateTestWindow() {
	className := "testClass"

	instance := GetModuleHandle()

	cursor, err := LoadCursorResource(IDC_ARROW)
	if err != nil {
		log.Println(err)
		return
	}

	wc := WNDCLASSEXW{
		LpfnWndProc:   syscall.NewCallback(testWindowCallback),
		HInstance:     instance,
		HCursor:       cursor,
		HbrBackground: COLOR_WINDOW + 1,
		LpszClassName: syscall.StringToUTF16Ptr(className),
	}
	wc.CbSize = UINT(unsafe.Sizeof(wc))

	if _, err = RegisterClassEx(&wc); err != nil {
		log.Println(err)
		return
	}

	_, err = CreateWindow(
		className,
		"Test Window",
		WS_VISIBLE | WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		0,
		0,
		instance,
	)
	if err != nil {
		log.Println(err)
		return
	}
}
