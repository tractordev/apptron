package win32

type ATOM uint16
type UINT uint32
type BOOL int32
type UINT_PTR uintptr
type LONG_PTR uintptr
type LARGE_INTEGER int64
type LONG int32
type WORD uint16
type DWORD uint32
type LPCWSTR *uint16
type CHAR  uint8
type WCHAR uint16

type HANDLE uintptr
type HWND HANDLE
type HMODULE HANDLE
type HICON HANDLE
type HCURSOR HICON
type HBRUSH HANDLE
type HINSTANCE HANDLE
type HMENU HANDLE

type WPARAM UINT_PTR
type LPARAM LONG_PTR
type LRESULT LONG_PTR
type WNDPROC func(hwnd HWND, msg uint32, wparam WPARAM, lparam LPARAM) LRESULT

// https://github.com/AllenDang/w32/blob/ad0a36d80adcd081d5c0dded8e97a009b486d1db/constants.go

const (
	NULL = 0
	TRUE  = 1
	FALSE = 0
)

const (
	SW_SHOW        = 5
)

const (
	PM_REMOVE = 0x0001
)

const (
	CW_USEDEFAULT = ^0x7fffffff
)

const (
	WS_VISIBLE      = 0x10000000

	WS_CAPTION     = 0x00C00000
	WS_MAXIMIZEBOX = 0x00010000
	WS_MINIMIZEBOX = 0x00020000
	WS_OVERLAPPED  = 0x00000000
	WS_SYSMENU     = 0x00080000
	WS_THICKFRAME  = 0x00040000

	WS_OVERLAPPEDWINDOW = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
)

const (
	WM_QUIT    = 18
	WM_DESTROY = 0x0002
	WM_CLOSE   = 0x0010
)

const (
	COLOR_WINDOW = 5
)

const (
	IDC_ARROW = 32512
)

// https://docs.microsoft.com/en-us/windows/win32/hidpi/dpi-awareness-context
type  DPI_AWARENESS_CONTEXT HANDLE
const DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = (DPI_AWARENESS_CONTEXT)(0xffffffff - 4)

// https://docs.microsoft.com/en-us/windows/win32/api/windef/ns-windef-point
type POINT struct {
	x LONG
	y LONG
}

// https://docs.microsoft.com/en-us/windows/win32/api/windef/ns-windef-rect
type RECT struct {
	left   LONG
	top    LONG
	right  LONG
	bottom LONG
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-msg
type MSG struct {
	hwnd     HWND
	message  UINT
	wParam   WPARAM
	lParam   LPARAM
	time     DWORD
	pt       POINT
	lPrivate DWORD
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-wndclassexw
type WNDCLASSEXW struct {
	cbSize        UINT
	style         UINT
	lpfnWndProc   uintptr // OR WNDPROC?
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     HINSTANCE
	hIcon         HICON
	hCursor       HCURSOR
	hbrBackground HBRUSH
	lpszMenuName  LPCWSTR
	lpszClassName LPCWSTR
	hIconSm       HICON
}

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

type NOTIFYICONDATAW struct {
	cbSize           DWORD
	hWnd             HWND
	uID              UINT
	uFlags           UINT
	uCallbackMessage UINT
	hIcon            HICON
	szTip            [128]WCHAR
	dwState          DWORD
	dwStateMask      DWORD
	szInfo           [256]WCHAR
	uTimeout         UINT
	uVersion         UINT
	szInfoTitle      [64]WCHAR
	dwInfoFlags      DWORD
	guidItem         GUID
	hBalloonIcon     HICON
}

type NOTIFYICONDATA NOTIFYICONDATAW