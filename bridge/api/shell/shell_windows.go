package shell

import (
	"tractor.dev/apptron/bridge/platform/win32"
)

func ShowNotification(n Notification) {
}

func ShowMessage(msg MessageDialog) bool {
	var flags win32.UINT = 0

	switch msg.Level {
	case "error":
		flags |= win32.MB_ICONERROR
	case "warning":
		flags |= win32.MB_ICONWARNING
	default:
		flags |= win32.MB_ICONINFORMATION
	}

	switch msg.Buttons {
	case "okcancel":
		flags |= win32.MB_OKCANCEL
	case "yesno":
		flags |= win32.MB_YESNO
	default:
		flags |= win32.MB_OK
	}

	return win32.MessageBox(win32.NULL, msg.Body, msg.Title, flags) == win32.IDOK
}

func ShowFilePicker(fd FileDialog) []string {
	return []string{}
}

func ReadClipboard() string {
	return win32.OS_GetClipboardText()
}

func WriteClipboard(text string) bool {
	return win32.OS_SetClipboardText(text)
}
