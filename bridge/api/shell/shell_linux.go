package shell

import "tractor.dev/apptron/bridge/platform/linux"

func ShowNotification(n Notification) {
}

func ShowMessage(msg MessageDialog) bool {
  return false
}

func ShowFilePicker(fd FileDialog) []string {
  return []string{}
}

func ReadClipboard() string {
  return linux.OS_GetClipboardText()
}

func WriteClipboard(text string) bool {
  return linux.OS_SetClipboardText(text)
}
