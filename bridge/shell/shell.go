package shell

/*
#include "../../lib/hostbridge.h"
*/
import "C"

type Notification struct {
	Title    string
	Subtitle string // for MacOS only
	Body     string
	/*
		Silent   bool
	*/
}

type FileDialog struct {
	Title     string
	Directory string
	Filename  string
	Mode      string   // pickfile, pickfiles, pickfolder, savefile
	Filters   []string // each string is comma delimited (go,rs,toml) with optional label prefix (text:go,txt)
}

type MessageDialog struct {
	Title   string
	Body    string
	Level   string // info, warning, error
	Buttons string // ok, okcancel, yesno
}

func ShowNotification(n Notification) {
	C.shell_show_notification(C.CString(n.Title), C.CString(n.Subtitle), C.CString(n.Body))
}

func ShowFilePicker(picker FileDialog) string {
	return ""
}

func ShowMessage(msg MessageDialog) bool {
	return false
}
