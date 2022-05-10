package win32

type ATOM uint16
type UINT uint32
type BOOL int32
type BYTE uint8
type UINT_PTR uintptr
type LONG_PTR uintptr
type ULONG_PTR uintptr
type LONG int32
type WORD uint16
type DWORD uint32
type LPCWSTR *uint16
type LPWSTR *uint16
type CHAR uint8
type WCHAR uint16

type LARGE_INTEGER struct {
	LowPart  DWORD
	HighPart DWORD
}

type HANDLE uintptr
type HWND HANDLE
type HMODULE HANDLE
type HICON HANDLE
type HCURSOR HICON
type HBRUSH HANDLE
type HINSTANCE HANDLE
type HMENU HANDLE
type HBITMAP HANDLE

type WPARAM UINT_PTR
type LPARAM LONG_PTR
type LRESULT LONG_PTR
type WNDPROC func(hwnd HWND, msg uint32, wparam WPARAM, lparam LPARAM) LRESULT

// https://github.com/AllenDang/w32/blob/ad0a36d80adcd081d5c0dded8e97a009b486d1db/constants.go

const (
	NULL  = 0
	TRUE  = 1
	FALSE = 0
)

const (
	SW_SHOW = 5
)

const (
	PM_REMOVE = 0x0001
)

const (
	CW_USEDEFAULT = ^0x7fffffff
)

const (
	WS_VISIBLE = 0x10000000

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
	WM_USER    = 0x0400

	WM_LBUTTONDOWN = 0x0201
	WM_RBUTTONDOWN = 0x0204
)

const (
	COLOR_WINDOW = 5
)

const (
	IDC_ARROW = 32512
)

const (
	NIF_MESSAGE = 0x00000001
	NIF_ICON    = 0x00000002
	NIF_TIP     = 0x00000004
	NIF_STATE   = 0x00000008
	NIF_INFO    = 0x00000010
)

const (
	NIM_ADD    = 0x00000000
	NIM_MODIFY = 0x00000001
	NIM_DELETE = 0x00000002
)

const (
	MIIM_BITMAP     = 0x00000080
	MIIM_CHECKMARKS = 0x00000008
	MIIM_DATA       = 0x00000020
	MIIM_FTYPE      = 0x00000100
	MIIM_ID         = 0x00000002
	MIIM_STATE      = 0x00000001
	MIIM_STRING     = 0x00000040
	MIIM_SUBMENU    = 0x00000004
	MIIM_TYPE       = 0x00000010

	MFT_STRING     = 0x00000000
	MFT_RADIOCHECK = 0x00000200
	MFT_SEPARATOR  = 0x00000800

	MFS_CHECKED   = 0x00000008
	MFS_DISABLED  = 0x00000003
	MFS_ENABLED   = 0x00000000
	MFS_UNCHECKED = 0x00000000

	TPM_CENTERALIGN = 0x0004
	TPM_LEFTALIGN   = 0x0000
	TPM_RIGHTALIGN  = 0x0008

	TPM_BOTTOMALIGN  = 0x0020
	TPM_TOPALIGN     = 0x0000
	TPM_VCENTERALIGN = 0x0010

	TPM_NONOTIFY  = 0x0080
	TPM_RETURNCMD = 0x0100

	TPM_LEFTBUTTON  = 0x0000
	TPM_RIGHTBUTTON = 0x0002
)

const UINT_MAX = ^uint(0)

// https://docs.microsoft.com/en-us/windows/win32/hidpi/dpi-awareness-context
const DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = (HANDLE)(UINT_MAX - 4 + 1)

// https://docs.microsoft.com/en-us/windows/win32/api/windef/ns-windef-point
type POINT struct {
	X LONG
	Y LONG
}

// https://docs.microsoft.com/en-us/windows/win32/api/windef/ns-windef-rect
type RECT struct {
	Left   LONG
	Top    LONG
	Right  LONG
	Bottom LONG
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-msg
type MSG struct {
	Hwnd     HWND
	Message  UINT
	WParam   WPARAM
	LParam   LPARAM
	Time     DWORD
	Pt       POINT
	LPrivate DWORD
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-wndclassexw
type WNDCLASSEXW struct {
	CbSize        UINT
	Style         UINT
	LpfnWndProc   uintptr // OR WNDPROC?
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     HINSTANCE
	HIcon         HICON
	HCursor       HCURSOR
	HbrBackground HBRUSH
	LpszMenuName  LPCWSTR
	LpszClassName LPCWSTR
	LIconSm       HICON
}

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

// https://docs.microsoft.com/en-us/windows/win32/api/shellapi/ns-shellapi-notifyicondataa
type NOTIFYICONDATAW struct {
	CbSize           DWORD
	HWnd             HWND
	UID              UINT
	UFlags           UINT
	UCallbackMessage UINT
	HIcon            HICON
	SzTip            [128]WCHAR // @Incomplete: according to the docs, sometimes this is 64 instead?
	DwState          DWORD
	DwStateMask      DWORD
	SzInfo           [256]WCHAR
	UTimeout         UINT
	UVersion         UINT
	SzInfoTitle      [64]WCHAR
	DwInfoFlags      DWORD
	GuidItem         GUID
	HBalloonIcon     HICON
}

type NOTIFYICONDATA NOTIFYICONDATAW

type MENUITEMINFOW struct {
	CbSize        UINT
	FMask         UINT
	FType         UINT
	FState        UINT
	WID           UINT
	HSubMenu      HMENU
	HbmpChecked   HBITMAP
	HbmpUnchecked HBITMAP
	DwItemData    ULONG_PTR
	DwTypeData    LPWSTR
	Cch           UINT
	HbmpItem      HBITMAP
}

type MENUITEMINFO MENUITEMINFOW
