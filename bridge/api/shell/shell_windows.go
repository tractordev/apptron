package shell

import (
	"tractor.dev/apptron/bridge/platform/win32"
)

func ShowNotification(n Notification) {
}

func ShowMessage(msg MessageDialog) bool {
	return false
}

func ShowFilePicker(fd FileDialog) []string {
	return []string{}
}

func ReadClipboard() string {
	return win32.OS_GetClipboardText()
}

func WriteClipboard(text string) bool {
	return false
}
