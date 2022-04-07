package shell

import (
	"strings"

	"github.com/progrium/macdriver/cocoa"
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
	alert := objc.Get("NSAlert").Alloc().Init()

	switch msg.Level {
	case "error":
		alert.Set("alertStyle:", objc.Get("NSAlertStyleCritical"))
	case "warning":
		alert.Set("alertStyle:", objc.Get("NSAlertStyleWarning"))
	default:
		alert.Set("alertStyle:", objc.Get("NSAlertStyleInformational"))
	}

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

	alert.Set("messageText:", core.String(msg.Title))
	alert.Set("informativeText:", core.String(msg.Body))

	ret := alert.Send("runModal")
	return ret.Int() == 1000
}

func ShowFilePicker(fd FileDialog) []string {
	var picker objc.Object
	switch fd.Mode {
	case "savefile":
		picker = objc.Get("NSSavePanel").Send("savePanel")
	case "pickfiles":
		picker = objc.Get("NSOpenPanel").Send("openPanel")
		picker.Set("allowsMultipleSelection:", core.True)
	case "pickfolder":
		picker = objc.Get("NSOpenPanel").Send("openPanel")
		picker.Set("canChooseDirectories:", core.True)
	default: // pickfile
		picker = objc.Get("NSOpenPanel").Send("openPanel")

	}
	if fd.Filename != "" {
		picker.Set("nameFieldStringValue:", core.String(fd.Filename))
	}
	if fd.Directory != "" {
		url := core.NSURL_fileURLWithPath_isDirectory_(core.String(fd.Directory), true)
		picker.Set("directoryURL:", url)
	}
	if fd.Filters != nil {
		var filters []objc.Object
		for _, entry := range fd.Filters {
			kvp := strings.Split(entry, ":")
			var idx int
			if len(kvp) > 1 {
				idx = 1
			}
			for _, ext := range strings.Split(kvp[idx], ",") {
				filters = append(filters, core.String(ext))
			}
		}
		picker.Set("allowedFileTypes:", core.NSArray_WithObjects(filters...))
	}
	picker.Set("title:", core.String(fd.Title))
	ret := picker.Send("runModal")
	if ret.Int() == 1 {
		urls := core.NSArray_fromRef(picker.Send("URLs"))
		count := int(urls.Count())
		paths := make([]string, count)
		for i := 0; i < count; i++ {
			o := urls.ObjectAtIndex(uint64(i))
			paths[i] = o.Get("path").String()
		}
		return paths
	}
	return []string{}
}

func ReadClipboard() string {
	pb := cocoa.NSPasteboard_GeneralPasteboard()
	return pb.StringForType(cocoa.NSPasteboardTypeString)
}

func WriteClipboard(text string) bool {
	pb := cocoa.NSPasteboard_GeneralPasteboard()
	pb.ClearContents()
	pb.SetStringForType(text, cocoa.NSPasteboardTypeString)
	return true
}
