package shell

/*
#include "../../lib/hostbridge.h"
*/
import "C"

import (
	"strings"
	"unsafe"
)

var Module *module

func init() {
	Module = &module{}
}

type module struct{}

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

func (m module) ShowNotification(n Notification) {
	ShowNotification(n)
}

func ShowMessage(msg MessageDialog) bool {
	return Module.ShowMessage(msg)
}

func (m module) ShowMessage(msg MessageDialog) bool {
	return fromCBool(
		C.shell_show_dialog(C.CString(msg.Title), C.CString(msg.Body), C.CString(msg.Level), C.CString(msg.Buttons)))
}

func ShowFilePicker(fd FileDialog) []string {
	return Module.ShowFilePicker(fd)
}

func (m module) ShowFilePicker(fd FileDialog) []string {
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
	return Module.ReadClipboard()
}

func (m module) ReadClipboard() string {
	C.reset_temporary_storage()
	return C.GoString(C.shell_read_clipboard())
}

func WriteClipboard(text string) bool {
	return Module.WriteClipboard(text)
}

func (m module) WriteClipboard(text string) bool {
	return fromCBool(C.shell_write_clipboard(C.CString(text)))
}

func RegisterShortcut(accelerator string) bool {
	return Module.RegisterShortcut(accelerator)
}

func (m module) RegisterShortcut(accelerator string) bool {
	return fromCBool(C.shell_register_shortcut(C.CString(accelerator)))
}

func IsShortcutRegistered(accelerator string) bool {
	return Module.IsShortcutRegistered(accelerator)
}

func (m module) IsShortcutRegistered(accelerator string) bool {
	return fromCBool(C.shell_is_shortcut_registered(C.CString(accelerator)))
}

func UnregisterShortcut(accelerator string) bool {
	return Module.UnregisterShortcut(accelerator)
}

func (m module) UnregisterShortcut(accelerator string) bool {
	return fromCBool(C.shell_unregister_shortcut(C.CString(accelerator)))
}

func UnregisterAllShortcuts() bool {
	return Module.UnregisterAllShortcuts()
}

func (m module) UnregisterAllShortcuts() bool {
	return fromCBool(C.shell_unregister_all_shortcuts())
}

func fromCBool(it C.uchar) bool {
	return int(it) != 0
}
