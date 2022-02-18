package shell

/*
#include "../../lib/hostbridge.h"
*/
import "C"

import (
	"fmt"
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
	C.reset_temporary_storage()

	/*
	var filters C.StringArray
	filterCount := len(fd.Filters)
	if filterCount > 0 {
		cstr := make([]*C.char, filterCount)
		for i, it := range fd.Filters {
			cstr[i] = C.CString(it)
		}

		filters = C.StringArray{ data: (**C.char)(unsafe.Pointer(&cstr[0])), count: C.int(filterCount) }
	} else {
		filters = C.StringArray{ data: (**C.char)(nil), count: C.int(0) }
	}
	*/

	files := C.shell_show_file_picker(C.CString(fd.Title), C.CString(fd.Directory), C.CString(fd.Filename), C.CString(fd.Mode))

	n := int(files.count)
	result := make([]string, n)

	fileData := (*[1 << 28]*C.char)(unsafe.Pointer(files.data))[:n:n]
	for i := 0; i < n; i++ {
		str := C.GoString(fileData[i])
		fmt.Println("str", str)
		result[i] = str
	}

	return result
}

func toBool(it C.uchar) bool {
	return int(it) != 0
}
