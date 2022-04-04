package shell

import (
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

func ShowNotification(n Notification) {
	notification := objc.Get("NSUserNotification").Alloc().Init()
	notification.Set("title:", core.String(n.Title))
	notification.Set("informativeText:", core.String(n.Body))

	center := objc.Get("NSUserNotificationCenter").Send("defaultUserNotificationCenter")
	center.Send("deliverNotification:", notification)
	notification.Release()
}

func ShowMessage(msg MessageDialog) bool {
	return false
}

func ShowFilePicker(fd FileDialog) []string {
	return []string{}
}

func ReadClipboard() string {
	return ""
}

func WriteClipboard(text string) bool {
	return true
}

func RegisterShortcut(accelerator string) bool {
	return true
}

func IsShortcutRegistered(accelerator string) bool {
	return true
}

func UnregisterShortcut(accelerator string) bool {
	return true
}

func UnregisterAllShortcuts() bool {
	return true
}
