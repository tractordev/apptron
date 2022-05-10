package shell

func ShowNotification(n Notification) {
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
	return false
}
