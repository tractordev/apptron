package shell

func ShowNotification(n Notification) {
}

func ShowMessage(msg MessageDialog) bool {
	/*
	switch msg.Level {
	case "error":
		alert.Set("alertStyle:", objc.Get("NSAlertStyleCritical"))
	case "warning":
		alert.Set("alertStyle:", objc.Get("NSAlertStyleWarning"))
	default:
		alert.Set("alertStyle:", objc.Get("NSAlertStyleInformational"))
	}
	*/

	/*
	switch msg.Buttons {
	case "ok":
		alert.Send("addButtonWithTitle:", core.String("OK"))
	case "okcancel":
		alert.Send("addButtonWithTitle:", core.String("OK"))
		alert.Send("addButtonWithTitle:", core.String("Cancel"))
	case "yesno":
		alert.Send("addButtonWithTitle:", core.String("Yes"))
		alert.Send("addButtonWithTitle:", core.String("No"))
	}
	*/

	/*
	alert.Set("messageText:", core.String(msg.Title))
	alert.Set("informativeText:", core.String(msg.Body))
	*/

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
