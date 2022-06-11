//go:build windows

package win32

import (
	"log"
	"reflect"
	"syscall"
	"unsafe"
)

//
// Imports
//

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	pGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
	pExitProcess      = kernel32.NewProc("ExitProcess")

	pGlobalLock    = kernel32.NewProc("GlobalLock")
	pGlobalUnlock  = kernel32.NewProc("GlobalUnlock")
	pGlobalAlloc   = kernel32.NewProc("GlobalAlloc")
	pGlobalFree    = kernel32.NewProc("GlobalFree")
	pRtlMoveMemory = kernel32.NewProc("RtlMoveMemory")
)

func GetModuleHandle() HINSTANCE {
	ret, _, _ := pGetModuleHandleW.Call(uintptr(0))
	return HINSTANCE(ret)
}

func ExitProcess(exitCode UINT) {
	pExitProcess.Call(uintptr(exitCode))
}

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	pCreateWindowExW     = user32.NewProc("CreateWindowExW")
	pDefWindowProcW      = user32.NewProc("DefWindowProcW")
	pDestroyWindow       = user32.NewProc("DestroyWindow")
	pSetWindowPos        = user32.NewProc("SetWindowPos")
	pShowWindow          = user32.NewProc("ShowWindow")
	pGetWindowPlacement  = user32.NewProc("GetWindowPlacement")
	pSetWindowPlacement  = user32.NewProc("SetWindowPlacement")
	pMonitorFromWindow   = user32.NewProc("MonitorFromWindow")
	pSetWindowTextW      = user32.NewProc("SetWindowTextW")
	pGetCursorPos        = user32.NewProc("GetCursorPos")
	pSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	pGetActiveWindow     = user32.NewProc("GetActiveWindow")
	pSetWindowLongW      = user32.NewProc("SetWindowLongW")
	pGetWindowLongW      = user32.NewProc("GetWindowLongW")
	pValidateRect        = user32.NewProc("ValidateRect")
	pGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	pSetFocus            = user32.NewProc("SetFocus")

	pDispatchMessageW    = user32.NewProc("DispatchMessageW")
	pGetMessageW         = user32.NewProc("GetMessageW")
	pPeekMessageW        = user32.NewProc("PeekMessageW")
	pLoadCursorW         = user32.NewProc("LoadCursorW")
	pPostQuitMessage     = user32.NewProc("PostQuitMessage")
	pRegisterClassExW    = user32.NewProc("RegisterClassExW")
	pTranslateMessage    = user32.NewProc("TranslateMessage")
	pEnumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
	pEnumDisplaySettings = user32.NewProc("EnumDisplaySettingsW")
	pGetMonitorInfoW     = user32.NewProc("GetMonitorInfoW")

	pOpenClipboard    = user32.NewProc("OpenClipboard")
	pCloseClipboard   = user32.NewProc("CloseClipboard")
	pGetClipboardData = user32.NewProc("GetClipboardData")
	pEmptyClipboard   = user32.NewProc("EmptyClipboard")
	pSetClipboardData = user32.NewProc("SetClipboardData")

	pCreateMenu       = user32.NewProc("CreateMenu")
	pCreatePopupMenu  = user32.NewProc("CreatePopupMenu")
	pDestroyMenu      = user32.NewProc("DestroyMenu")
	pTrackPopupMenu   = user32.NewProc("TrackPopupMenu")
	pInsertMenuItemW  = user32.NewProc("InsertMenuItemW")
	pGetMenuItemCount = user32.NewProc("GetMenuItemCount")

	pCreateIconFromResourceEx    = user32.NewProc("CreateIconFromResourceEx")
	pLookupIconIdFromDirectoryEx = user32.NewProc("LookupIconIdFromDirectoryEx")

	pSetProcessDpiAwarenessContext = user32.NewProc("SetProcessDpiAwarenessContext")

	pMessageBoxW = user32.NewProc("MessageBoxW")
)

func CreateWindowExW(dwExStyle DWORD, className string, windowName string, style uint32, x, y, width, height int32, parent, menu, instance HINSTANCE, lpParam uintptr) (HWND, error) {
	ret, _, err := pCreateWindowExW.Call(
		uintptr(dwExStyle),
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
		uintptr(lpParam),
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

func DestroyWindow(hwnd HWND) bool {
	ret, _, _ := pDestroyWindow.Call(uintptr(hwnd))
	return int32(ret) != 0
}

func SetWindowPos(hwnd HWND, hwndInsertAfter HWND, x int, y int, cx int, cy int, flags UINT) bool {
	ret, _, _ := pSetWindowPos.Call(uintptr(hwnd), uintptr(hwndInsertAfter), uintptr(x), uintptr(y), uintptr(cx), uintptr(cy), uintptr(flags))
	return int32(ret) != 0
}

func MonitorFromWindow(hwnd HWND, dwFlags DWORD) HMONITOR {
	ret, _, _ := pMonitorFromWindow.Call(uintptr(hwnd), uintptr(dwFlags))
	return HMONITOR(ret)
}

func ShowWindow(hwnd HWND, nCmdShow int) bool {
	ret, _, _ := pShowWindow.Call(uintptr(hwnd), uintptr(nCmdShow))
	return int32(ret) != 0
}

func SetWindowTextW(hwnd HWND, title string) bool {
	ret, _, _ := pSetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))))
	return int32(ret) != 0
}

func GetClientRect(hwnd HWND, lpRect *RECT) bool {
	ret, _, _ := pShowWindow.Call(uintptr(hwnd), uintptr(unsafe.Pointer(lpRect)))
	return int32(ret) != 0
}

func GetForegroundWindow() HWND {
	ret, _, _ := pGetForegroundWindow.Call()
	return HWND(ret)
}

func SetFocus(hwnd HWND) HWND {
	ret, _, _ := pSetFocus.Call(uintptr(hwnd))
	return HWND(ret)
}

func ValidateRect(hwnd HWND, lpRect *RECT) bool {
	ret, _, _ := pValidateRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(lpRect)))
	return int32(ret) != 0
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

func RegisterClassExW(wcx *WNDCLASSEXW) (uint16, error) {
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

func SetWindowLongW(hwnd HWND, index int, long LONG) LONG {
	ret, _, _ := pSetWindowLongW.Call(uintptr(hwnd), uintptr(index), uintptr(long))
	return LONG(ret)
}

func GetWindowLongW(hwnd HWND, index int) LONG {
	ret, _, _ := pGetWindowLongW.Call(uintptr(hwnd), uintptr(index))
	return LONG(ret)
}

func EnumDisplayMonitors(hdc HDC, clip *RECT, enumProc MONITORENUMPROC, data LPARAM) bool {
	ret, _, _ := pEnumDisplayMonitors.Call(uintptr(hdc), uintptr(unsafe.Pointer(clip)), syscall.NewCallback(enumProc), uintptr(data))
	return ret != 0
}

func EnumDisplaySettings(deviceName *uint16, iModeNum DWORD, lpDevMode *DEVMODE) bool {
	lpDevMode.DmSize = WORD(unsafe.Sizeof(*lpDevMode))

	ret, _, _ := pEnumDisplaySettings.Call(uintptr(unsafe.Pointer(deviceName)), uintptr(iModeNum), uintptr(unsafe.Pointer(lpDevMode)))
	return ret != 0
}

func GetMonitorInfo(monitor HMONITOR, info *MONITORINFOEX) bool {
	info.CbSize = DWORD(unsafe.Sizeof(*info))

	ret, _, _ := pGetMonitorInfoW.Call(uintptr(monitor), uintptr(unsafe.Pointer(info)))
	return ret != 0
}

func GetCursorPos(pos *POINT) bool {
	ret, _, _ := pGetCursorPos.Call(uintptr(unsafe.Pointer(pos)))
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

func CreateIconFromResourceEx(bytes *BYTE, size DWORD, icon BOOL, ver DWORD, cxDesired int32, cyDesired int32, flags UINT) HICON {
	result, _, _ := pCreateIconFromResourceEx.Call(
		uintptr(unsafe.Pointer(bytes)),
		uintptr(size),
		uintptr(icon),
		uintptr(ver),
		uintptr(cxDesired),
		uintptr(cyDesired),
		uintptr(flags),
	)
	return HICON(result)
}

func LookupIconIdFromDirectoryEx(bytes *BYTE, icon BOOL, cxDesired int32, cyDesired int32, flags UINT) int32 {
	result, _, _ := pLookupIconIdFromDirectoryEx.Call(
		uintptr(unsafe.Pointer(bytes)),
		uintptr(icon),
		uintptr(cxDesired),
		uintptr(cyDesired),
		uintptr(flags),
	)
	return int32(result)
}

func MessageBox(hwnd HWND, text string, caption string, flags UINT) int32 {
	result, _, _ := pMessageBoxW.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		uintptr(flags),
	)
	return int32(result)
}

var (
	shell32 = syscall.NewLazyDLL("shell32.dll")

	pShell_NotifyIconW = shell32.NewProc("Shell_NotifyIconW")
)

func Shell_NotifyIconW(dwMessage DWORD, nid *NOTIFYICONDATA) bool {
	ret, _, _ := pShell_NotifyIconW.Call(uintptr(dwMessage), uintptr(unsafe.Pointer(nid)))
	return ret != 0
}

var (
	shcore = syscall.NewLazyDLL("shcore.dll")

	// min support Windows 8.1 [desktop apps only]
	pGetDpiForMonitor = shcore.NewProc("GetDpiForMonitor")
)

func GetDpiForMonitor(monitor HMONITOR, dpiType uint32 /*MONITOR_DPI_TYPE*/, dpiX *UINT, dpiY *UINT) bool {
	ret, _, _ := pGetDpiForMonitor.Call(uintptr(monitor), uintptr(dpiType), uintptr(unsafe.Pointer(dpiX)), uintptr(unsafe.Pointer(dpiY)))
	return ret == 0 /*S_OK*/
}

var (
	winmm = syscall.NewLazyDLL("winmm.dll")

	pTimeBeginPeriod = winmm.NewProc("timeBeginPeriod")
)

//
// Helpers
//

func MakeIntResource(id uint16) *uint16 {
	return (*uint16)(unsafe.Pointer(uintptr(id)))
}

func LOWORD(dw uint32) uint16 {
	return uint16(dw)
}

func HIWORD(dw uint32) uint16 {
	return uint16(dw >> 16 & 0xffff)
}

func Utf16PtrToString(p uintptr) string {
	n := 0
	for ptr := unsafe.Pointer(p); *(*uint16)(ptr) != 0; n++ {
		ptr = unsafe.Pointer(uintptr(ptr) +
			unsafe.Sizeof(*((*uint16)(unsafe.Pointer(p)))))
	}

	var s []uint16
	h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	h.Data = p
	h.Len = n
	h.Cap = n
	return syscall.UTF16ToString(s)
}

//
// Functions
//

var win32SleepIsGranular = false

func OS_Init() {
	// NOTE(nick): request high-precision timers
	result, _, _ := pTimeBeginPeriod.Call(1)
	win32SleepIsGranular = UINT(result) == 0 /* TIMERR_NOERROR */

	//log.Println("[OS] sleep is granular", win32SleepIsGranular)
}

/*
func SleepMS(float64 miliseconds) {
  // @Incomplete: only do this if sleep is granular!
  // Otherwise do some sort of busy wait thing

  LARGE_INTEGER ft;
  ft.QuadPart = -(10 * (__int64)(miliseconds * 1000));

  timer := CreateWaitableTimer(NULL, TRUE, NULL);
  SetWaitableTimer(timer, &ft, 0, NULL, NULL, 0);
  WaitForSingleObject(timer, INFINITE);
  CloseHandle(timer);
}
*/

func OS_GetClipboardText() string {
	var result string

	ret, _, _ := pOpenClipboard.Call(uintptr(NULL))
	if ret == 0 {
		log.Println("[clipboard] Failed to open clipboard.")
		return result
	}

	ret, _, _ = pGetClipboardData.Call(uintptr(CF_UNICODETEXT))
	handle := HANDLE(ret)
	if handle == 0 {
		log.Println("[clipboard] Failed to convert clipboard to string.")
		pCloseClipboard.Call()
		return result
	}

	ret, _, _ = pGlobalLock.Call(uintptr(handle))
	defer pGlobalUnlock.Call(uintptr(handle))

	if ret == 0 {
		log.Println("[clipboard] Failed to lock global handle.")
		pCloseClipboard.Call()
		return result
	}

	result = Utf16PtrToString(ret)

	pCloseClipboard.Call()

	return result
}

func OS_SetClipboardText(text string) bool {
	s, err := syscall.UTF16FromString(text)
	if err != nil {
		log.Println("[clipboard] Failed to convert string to utf16: %w", err)
		return false
	}

	hMem, _, err := pGlobalAlloc.Call(GMEM_MOVEABLE, uintptr(len(s)*int(unsafe.Sizeof(s[0]))))
	if hMem == 0 {
		log.Println("[clipboard] Failed to alloc global memory: %w", err)
		return false
	}

	p, _, err := pGlobalLock.Call(hMem)
	if p == 0 {
		log.Println("[clipboard] Failed to lock global memory: %w", err)
		return false
	}
	defer pGlobalUnlock.Call(hMem)

	pRtlMoveMemory.Call(p, uintptr(unsafe.Pointer(&s[0])), uintptr(len(s)*int(unsafe.Sizeof(s[0]))))

	ret, _, _ := pOpenClipboard.Call(uintptr(NULL))
	if ret == 0 {
		log.Println("[clipboard] Failed to open clipboard.")
		return false
	}
	defer pCloseClipboard.Call()

	r, _, err := pEmptyClipboard.Call()
	if r == 0 {
		log.Println("[clipboard] Failed to clear clipboard: %w", err)
		return false
	}

	v, _, err := pSetClipboardData.Call(CF_UNICODETEXT, hMem)
	if v == 0 {
		pGlobalFree.Call(hMem)
		log.Println("[clipboard] Failed to set text to clipboard: %w", err)
		return false
	}

	return true
}

func PollEvents() {
	for {
		msg := MSG{}
		if !PeekMessageW(&msg, NULL, 0, 0, PM_REMOVE) {
			break
		}

		switch msg.Message {
		case WM_QUIT:
			PostQuitMessage(0)
		}

		TranslateMessage(&msg)
		DispatchMessage(&msg)
	}
}

func (info *MONITORINFOEX) GetDeviceName() string {
	return syscall.UTF16ToString(info.DeviceName[:])
}

func MakeMenuItemSeparator() MENUITEMINFO {
	result := MENUITEMINFO{}
	result.CbSize = UINT(unsafe.Sizeof(result))
	result.FMask = MIIM_ID | MIIM_DATA | MIIM_TYPE
	result.WID = 0
	result.DwItemData = 0
	return result
}

func MakeMenuItem(id int, label string, disabled bool, checked bool, isRadio bool) MENUITEMINFO {
	result := MENUITEMINFO{}

	result.CbSize = UINT(unsafe.Sizeof(result))
	result.FMask = MIIM_ID | MIIM_STATE | MIIM_DATA | MIIM_TYPE
	result.FType = MFT_STRING

	result.FState = 0
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

	result.WID = UINT(id)
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
var didInitTrayWindowClass = false

type Win32_Tray struct {
	iconData NOTIFYICONDATA
	window   HWND
	menu     HMENU
	callback func(id int32)
}

var trays = []Win32_Tray{}

const Win32TrayIconMessage = (WM_USER + 1)

func trayWindowCallback(hwnd HWND, message uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	switch message {
	case Win32TrayIconMessage:
		switch lParam {
		case WM_LBUTTONDOWN, WM_RBUTTONDOWN:

			SetForegroundWindow(hwnd)

			index := GetWindowLongW(hwnd, GWL_USERDATA)
			tray := trays[index]

			mousePosition := POINT{}
			GetCursorPos(&mousePosition)

			result := TrackPopupMenu(tray.menu, TPM_RIGHTBUTTON|TPM_NONOTIFY|TPM_RETURNCMD, int32(mousePosition.X), int32(mousePosition.Y), 0, hwnd, nil)

			if result > 0 {
				if tray.callback != nil {
					tray.callback(result)
				}
			}

		default:
			break
		}

	default:
		return DefWindowProc(hwnd, message, wParam, lParam)
	}
	return 0
}

func RegisterWindowClass(className string, instance HINSTANCE, callback WNDPROC, style UINT) bool {
	cursor, err := LoadCursorResource(IDC_ARROW)
	if err != nil {
		log.Println(err)
		return false
	}

	wc := WNDCLASSEXW{
		LpfnWndProc:   syscall.NewCallback(callback),
		HInstance:     instance,
		HCursor:       cursor,
		Style:         style,
		HbrBackground: COLOR_WINDOW + 1,
		LpszClassName: syscall.StringToUTF16Ptr(className),
	}
	wc.CbSize = UINT(unsafe.Sizeof(wc))

	if _, err = RegisterClassExW(&wc); err != nil {
		log.Println(err)
		return false
	}

	return true
}

func NewTrayMenu(menu HMENU, icon []byte, callback func(id int32)) bool {
	trayClassName := "APPTRON_TRAY_WINDOW_CLASS"

	if !didInitTrayWindowClass {
		if !RegisterWindowClass(trayClassName, GetModuleHandle(), trayWindowCallback, 0) {
			log.Println("Failed to register tray window class!")
			return false
		}

		didInitTrayWindowClass = true
	}

	hwnd, err := CreateWindowExW(0, trayClassName, "Tray Window", 0, 0, 0, 1, 1, 0, 0, GetModuleHandle(), 0)
	if err != nil {
		log.Println("Failed to create tray window!", err)
		return false
	}

	trayIconData := NOTIFYICONDATA{}
	trayIconData.CbSize = DWORD(unsafe.Sizeof(trayIconData))
	trayIconData.HWnd = hwnd
	trayIconData.UID = 0
	trayIconData.UFlags = NIF_MESSAGE | NIF_ICON | NIF_TIP
	trayIconData.UCallbackMessage = Win32TrayIconMessage

	// @Robustness: convert from PNG to ICO

	iconSize := len(icon)
	if iconSize > 0 {
		data := (*BYTE)(unsafe.Pointer(&icon[0]))

		offset := LookupIconIdFromDirectoryEx(data, TRUE, 0, 0, 0x00008000 /*LR_SHARED*/)

		if offset > 0 {
			data = (*BYTE)(unsafe.Pointer(&icon[offset]))
			trayIconData.HIcon = CreateIconFromResourceEx(data, DWORD(iconSize), TRUE, 0x00030000, 32, 32, 0 /*LR_DEFAULTCOLOR*/)
		}
	}

	// @Robustness: provide a default placeholder icon?
	//trayIconData.HIcon = LoadIcon(GetModuleHandle(0), MAKEINTRESOURCE(101));

	trayIconData.SzTip[0] = 0 // @Incomplete: we should put the app name here

	Shell_NotifyIconW(NIM_ADD, &trayIconData)

	tray := Win32_Tray{}
	tray.menu = menu
	tray.window = hwnd
	tray.iconData = trayIconData
	tray.callback = callback

	index := len(trays)
	SetWindowLongW(hwnd, GWL_USERDATA, LONG(index))

	trays = append(trays, tray)

	return true
}

func RemoveAllTrayMenus() {
	for _, it := range trays {
		Shell_NotifyIconW(NIM_DELETE, &it.iconData)
	}

	trays = make([]Win32_Tray, 0)
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

	if _, err = RegisterClassExW(&wc); err != nil {
		log.Println(err)
		return
	}

	_, err = CreateWindowExW(
		0,
		className,
		"Test Window",
		WS_VISIBLE|WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		0,
		0,
		instance,
		0,
	)
	if err != nil {
		log.Println(err)
		return
	}
}
