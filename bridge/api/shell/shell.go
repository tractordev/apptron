package shell

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

func (m module) ShowNotification(n Notification) {
	ShowNotification(n)
}

func (m module) ShowMessage(msg MessageDialog) bool {
	return ShowMessage(msg)
}

func (m module) ShowFilePicker(fd FileDialog) []string {
	return ShowFilePicker(fd)
}

func (m module) ReadClipboard() string {
	return ReadClipboard()
}

func (m module) WriteClipboard(text string) bool {
	return WriteClipboard(text)
}

func (m module) RegisterShortcut(accelerator string) {
	// hotkey does its own dispatch so
	// this avoids a deadlock assuming
	// all rpc calls are dispatched
	//go RegisterShortcut(accelerator)
}

func (m module) IsShortcutRegistered(accelerator string) bool {
	return IsShortcutRegistered(accelerator)
}

func (m module) UnregisterShortcut(accelerator string) bool {
	return UnregisterShortcut(accelerator)
}

func (m module) UnregisterAllShortcuts() {
	UnregisterAllShortcuts()
}
