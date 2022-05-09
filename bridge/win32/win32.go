package win32

import (
	"log"
	"syscall"
	"unsafe"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	pGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	pCreateWindowExW  = user32.NewProc("CreateWindowExW")
	pDefWindowProcW   = user32.NewProc("DefWindowProcW")
	pDestroyWindow    = user32.NewProc("DestroyWindow")
	pDispatchMessageW = user32.NewProc("DispatchMessageW")
	pGetMessageW      = user32.NewProc("GetMessageW")
	pPeekMessageW     = user32.NewProc("PeekMessageW")
	pLoadCursorW      = user32.NewProc("LoadCursorW")
	pPostQuitMessage  = user32.NewProc("PostQuitMessage")
	pRegisterClassExW = user32.NewProc("RegisterClassExW")
	pTranslateMessage = user32.NewProc("TranslateMessage")

	pSetProcessDpiAwarenessContext = user32.NewProc("SetProcessDpiAwarenessContext")
)

var (
	shell32 = syscall.NewLazyDLL("shell32.dll")

	pShell_NotifyIconW = shell32.NewProc("Shell_NotifyIconW")
)

func GetModuleHandle() HINSTANCE {
	ret, _, _ := pGetModuleHandleW.Call(uintptr(0))
	/*
	if ret == 0 {
		return 0, err
	}
	*/
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

func SetProcessDpiAwarenessContext(context DPI_AWARENESS_CONTEXT) bool {
	ret, _, _ := pSetProcessDpiAwarenessContext.Call(uintptr(context))
	return ret != 0
}

func Shell_NotifyIconW(dwMessage DWORD, nid *NOTIFYICONDATA) bool {
	ret, _, _ := pShell_NotifyIconW.Call(uintptr(dwMessage), uintptr(unsafe.Pointer(nid)))
	return ret != 0
}


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


func RegisterWindowClass(className string, instance HINSTANCE, callback WNDPROC) {
	cursor, err := LoadCursorResource(IDC_ARROW)
	if err != nil {
		log.Println(err)
		return
	}

	wc := WNDCLASSEXW{
		lpfnWndProc:   syscall.NewCallback(callback),
		hInstance:     instance,
		hCursor:       cursor,
		hbrBackground: COLOR_WINDOW + 1,
		lpszClassName: syscall.StringToUTF16Ptr(className),
	}
	wc.cbSize = UINT(unsafe.Sizeof(wc))

	if _, err = RegisterClassEx(&wc); err != nil {
		log.Println(err)
		return
	}
}

func trayWindowCallback(hwnd HWND, message uint32, wParam WPARAM, lParam LPARAM) LRESULT {
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

var trayIconData NOTIFYICONDATA

const Win32TrayIconMessage = (WM_USER + 1)

func SetupTray() {
  trayClassName := "trayTestClass"
  RegisterWindowClass(trayClassName, GetModuleHandle(), trayWindowCallback)

  trayWindow, err := CreateWindow(trayClassName, "Tray Window", 0, 0, 0, 1, 1, 0, 0, GetModuleHandle());
  if err != nil {
  	log.Println(err)
  	return
  }

	trayIconData = NOTIFYICONDATA{}
  trayIconData.cbSize = DWORD(unsafe.Sizeof(trayIconData))
  trayIconData.hWnd = trayWindow
  trayIconData.uID = 0
  trayIconData.uFlags = NIF_MESSAGE | NIF_ICON | NIF_TIP
  trayIconData.uCallbackMessage = Win32TrayIconMessage
  //trayIconData.hIcon = LoadIcon(GetModuleHandle(0), MAKEINTRESOURCE(101));
  trayIconData.szTip[0] = 0

  Shell_NotifyIconW(NIM_ADD, &trayIconData)

  msg := MSG{}
  for {
  	if (GetMessage(&msg, 0, 0, 0)) {
  		TranslateMessage(&msg)
  		DispatchMessage(&msg)
  	} else {
  		break
  	}
  }
}


func Main() {
	className := "testClass"

	instance := GetModuleHandle()
	cursor, err := LoadCursorResource(IDC_ARROW)
	if err != nil {
		log.Println(err)
		return
	}

	wc := WNDCLASSEXW{
		lpfnWndProc:   syscall.NewCallback(trayWindowCallback),
		hInstance:     instance,
		hCursor:       cursor,
		hbrBackground: COLOR_WINDOW + 1,
		lpszClassName: syscall.StringToUTF16Ptr(className),
	}
	wc.cbSize = UINT(unsafe.Sizeof(wc))

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
