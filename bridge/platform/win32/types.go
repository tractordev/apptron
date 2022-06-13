//go:build windows

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
type HDC HANDLE
type HMONITOR HANDLE

type WPARAM UINT_PTR
type LPARAM LONG_PTR
type LRESULT LONG_PTR
type WNDPROC func(hwnd HWND, msg uint32, wparam WPARAM, lparam LPARAM) LRESULT

type MONITORENUMPROC func(unnamedParam1 HMONITOR, unnamedParam2 HDC, unnamedParam3 *RECT, unnamedParam4 LPARAM) uintptr

// https://github.com/AllenDang/w32/blob/ad0a36d80adcd081d5c0dded8e97a009b486d1db/constants.go

const UINT_MAX = ^uint(0)
const INT_MAX = ^int(0)

const LONG_MAX = 2147483647

const (
	NULL  = 0
	TRUE  = 1
	FALSE = 0
)

const (
	SW_HIDE            = 0
	SW_NORMAL          = 1
	SW_SHOWMINIMIZED   = 2
	SW_SHOWMAXIMIZED   = 3
	SW_MAXIMIZE        = 3
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_RESTORE         = 9
)

const (
	CS_VREDRAW         = 0x00000001
	CS_HREDRAW         = 0x00000002
	CS_KEYCVTWINDOW    = 0x00000004
	CS_DBLCLKS         = 0x00000008
	CS_OWNDC           = 0x00000020
	CS_CLASSDC         = 0x00000040
	CS_PARENTDC        = 0x00000080
	CS_NOKEYCVT        = 0x00000100
	CS_NOCLOSE         = 0x00000200
	CS_SAVEBITS        = 0x00000800
	CS_BYTEALIGNCLIENT = 0x00001000
	CS_BYTEALIGNWINDOW = 0x00002000
	CS_GLOBALCLASS     = 0x00004000
	CS_IME             = 0x00010000
	CS_DROPSHADOW      = 0x00020000
)

const (
	HWND_TOP       = (HWND)(0)
	HWND_BOTTOM    = (HWND)(1)
	HWND_TOPMOST   = (HWND)(UINT_MAX - 1 + 1)
	HWND_NOTOPMOST = (HWND)(UINT_MAX - 2 + 1)
)

const (
	MONITOR_DEFAULTTOPRIMARY = 0x00000001
	MONITOR_DEFAULTTONEAREST = 0x00000002
)

const (
	SWP_NOSIZE        = 0x0001
	SWP_NOMOVE        = 0x0002
	SWP_NOZORDER      = 0x0004
	SWP_NOACTIVATE    = 0x0010
	SWP_FRAMECHANGED  = 0x0020
	SWP_NOOWNERZORDER = 0x0200
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
	WS_POPUP       = 0x80000000

	WS_OVERLAPPEDWINDOW = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
)

const (
	WM_CREATE           = 0x0001
	WM_DESTROY          = 0x0002
	WM_MOVE             = 0x0003
	WM_SIZE             = 0x0005
	WM_ACTIVATE         = 0x0006
	WM_SETFOCUS         = 0x0007
	WM_KILLFOCUS        = 0x0008
	WM_CLOSE            = 0x0010
	WM_QUIT             = 0x0012
	WM_GETMINMAXINFO    = 0x0024
	WM_WINDOWPOSCHANGED = 0x0047
	WM_CHAR             = 0x0102
	WM_SYSCHAR          = 0x0106
	WM_UNICHAR          = 0x0109
	WM_SYSCOMMAND       = 0x0112
	WM_LBUTTONDOWN      = 0x0201
	WM_RBUTTONDOWN      = 0x0204
	WM_MOVING           = 0x0216
	WM_DPICHANGED       = 0x02E0
	WM_USER             = 0x0400
)

const (
	SC_KEYMENU      = 0xF100
	SC_SCREENSAVE   = 0xF140
	SC_MONITORPOWER = 0xF170
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

// https://docs.microsoft.com/en-us/windows/win32/hidpi/dpi-awareness-context
const DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = (HANDLE)(UINT_MAX - 4 + 1)

const PROCESS_PER_MONITOR_DPI_AWARE = 2

const (
	GWL_STYLE     = (int)(INT_MAX - 16 + 1)
	GWL_USERDATA  = (int)(INT_MAX - 21 + 1)
	GWLP_USERDATA = (int)(INT_MAX - 21 + 1)
)

const ENUM_CURRENT_SETTINGS = 0xFFFFFFFF

const CCHDEVICENAME = 32
const CCHFORMNAME = 32

const USER_DEFAULT_SCREEN_DPI = 96

const (
	CF_TEXT        = 1
	CF_UNICODETEXT = 13
)

const (
	GMEM_MOVEABLE = 0x0002
)

const (
	MB_OK       = 0x00000000
	MB_OKCANCEL = 0x00000001
	MB_YESNO    = 0x00000004

	MB_ICONWARNING     = 0x00000030
	MB_ICONINFORMATION = 0x00000040
	MB_ICONERROR       = 0x00000010

	IDOK = 1
)

const (
	DWMWA_CLOAKED = 14
)

const (
	LOGPIXELSX = 88
	LOGPIXELSY = 90
)

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

type MONITORINFO struct {
	CbSize    DWORD
	RcMonitor RECT
	RcWork    RECT
	DwFlags   DWORD
}

type MONITORINFOEXW struct {
	MONITORINFO
	DeviceName [CCHDEVICENAME]uint16
}

type MONITORINFOEX MONITORINFOEXW

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd183565.aspx
type DEVMODE struct {
	DmDeviceName    [CCHDEVICENAME]uint16
	DmSpecVersion   WORD
	DmDriverVersion WORD
	DmSize          WORD
	DmDriverExtra   WORD
	DmFields        DWORD

	// union!
	DmPosition POINT // 64 bytes
	/*
		DmOrientation   int16
		DmPaperSize     int16
		DmPaperLength   int16
		DmPaperWidth    int16
	*/
	_DmScale         int16
	_DmCopies        int16
	_DmDefaultSource int16
	_DmPrintQuality  int16

	DmColor            int16
	DmDuplex           int16
	DmYResolution      int16
	DmTTOption         int16
	DmCollate          int16
	DmFormName         [CCHFORMNAME]WCHAR
	DmLogPixels        WORD
	DmBitsPerPel       DWORD
	DmPelsWidth        DWORD
	DmPelsHeight       DWORD
	DmDisplayFlags     DWORD
	DmDisplayFrequency DWORD
	DmICMMethod        DWORD
	DmICMIntent        DWORD
	DmMediaType        DWORD
	DmDitherType       DWORD
	DmReserved1        DWORD
	DmReserved2        DWORD
	DmPanningWidth     DWORD
	DmPanningHeight    DWORD
}

type WINDOWPLACEMENT struct {
	Length           UINT
	Flags            UINT
	ShowCmd          UINT
	PtMinPosition    POINT
	PtMaxPosition    POINT
	RcNormalPosition RECT
	RcDevice         RECT
}

type MINMAXINFO struct {
	PtReserved     POINT
	PtMaxSize      POINT
	PtMaxPosition  POINT
	PtMinTrackSize POINT
	PtMaxTrackSize POINT
}
