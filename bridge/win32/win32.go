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
	pPeekMessageW     = user32.NewProc("PeekMessageW")
	pLoadCursorW      = user32.NewProc("LoadCursorW")
	pPostQuitMessage  = user32.NewProc("PostQuitMessage")
	pRegisterClassExW = user32.NewProc("RegisterClassExW")
	pTranslateMessage = user32.NewProc("TranslateMessage")

	pSetProcessDpiAwarenessContext = user32.NewProc("SetProcessDpiAwarenessContext")
)

func GetModuleHandle() (HINSTANCE, error) {
	ret, _, err := pGetModuleHandleW.Call(uintptr(0))
	if ret == 0 {
		return 0, err
	}
	return HINSTANCE(ret), nil
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

func SetProcessDpiAwarenessContext(context DPI_AWARENESS_CONTEXT) BOOL {
	ret, _, _ := pSetProcessDpiAwarenessContext.Call(uintptr(context))
	return BOOL(ret)
}

func trayWindowCallback(hwnd HWND, msg uint32, wparam WPARAM, lparam LPARAM) LRESULT {
	switch msg {
		case WM_CLOSE:
			DestroyWindow(hwnd)
		case WM_DESTROY:
			PostQuitMessage(0)
		default:
			ret := DefWindowProc(hwnd, msg, wparam, lparam)
			return ret
	}
	return 0
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

func Main() {
	className := "testClass"

	instance, err := GetModuleHandle()
	if err != nil {
		log.Println(err)
		return
	}

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
