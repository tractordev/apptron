package shell

/*
#include "../../lib/hostbridge.h"
*/
import "C"

import (
	"strings"
	"unsafe"
)

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

func ShowMessage(msg MessageDialog) bool {
	result := C.shell_show_dialog(C.CString(msg.Title), C.CString(msg.Body), C.CString(msg.Level), C.CString(msg.Buttons))
	return toBool(result)
}

func ShowFilePicker(fd FileDialog) []string {
	// @Cleanup: reset these at frame boundaries in the event loop?
	C.reset_temporary_storage()

	filters := strings.Join(fd.Filters, "|")

	files := C.shell_show_file_picker(C.CString(fd.Title), C.CString(fd.Directory), C.CString(fd.Filename), C.CString(fd.Mode), C.CString(filters))

	n := int(files.count)
	result := make([]string, n)

	fileData := (*[1 << 28]*C.char)(unsafe.Pointer(files.data))[:n:n]
	for i := 0; i < n; i++ {
		str := C.GoString(fileData[i])
		result[i] = str
	}

	return result
}

func ReadClipboard() string {
	C.reset_temporary_storage()

	result := C.shell_read_clipboard()
	return C.GoString(result)
}

func WriteClipboard(text string) bool {
	result := C.shell_write_clipboard(C.CString(text))
	return toBool(result)
}

func RegisterShortcut(accelerator string) bool {
	result := C.shell_register_shortcut(C.CString(accelerator))
	return toBool(result)
}

func IsShortcutRegistered(accelerator string) bool {
	result := C.shell_is_shortcut_registered(C.CString(accelerator))
	return toBool(result)
}

func UnregisterShortcut(accelerator string) bool {
	result := C.shell_unregister_shortcut(C.CString(accelerator))
	return toBool(result)
}

func UnregisterAllShortcuts() bool {
	result := C.shell_unregister_all_shortcuts()
	return toBool(result)
}

func toBool(it C.uchar) bool {
	return int(it) != 0
}
