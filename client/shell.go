package client

import (
	"context"

	"github.com/progrium/qtalk-go/fn"
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

type ShellModule struct {
	client *Client

	OnShortcut func(event Event)
}

// ShowNotification
func (m *ShellModule) ShowNotification(ctx context.Context, n Notification) (err error) {
	_, err = m.client.Call(ctx, "shell.ShowNotification", fn.Args{n}, nil)
	return
}

// ShowMessage
func (m *ShellModule) ShowMessage(ctx context.Context, msg MessageDialog) (ret bool, err error) {
	_, err = m.client.Call(ctx, "shell.ShowMessage", fn.Args{msg}, &ret)
	return
}

// ShowFilePicker
func (m *ShellModule) ShowFilePicker(ctx context.Context, fd FileDialog) (ret []string, err error) {
	_, err = m.client.Call(ctx, "shell.ShowFilePicker", fn.Args{fd}, &ret)
	return
}

// ReadClipboard
func (m *ShellModule) ReadClipboard(ctx context.Context) (ret string, err error) {
	_, err = m.client.Call(ctx, "shell.ReadClipboard", fn.Args{}, &ret)
	return
}

// WriteClipboard
func (m *ShellModule) WriteClipboard(ctx context.Context, text string) (ret bool, err error) {
	_, err = m.client.Call(ctx, "shell.WriteClipboard", fn.Args{text}, &ret)
	return
}

// RegisterShortcut
func (m *ShellModule) RegisterShortcut(ctx context.Context, accelerator string) (err error) {
	_, err = m.client.Call(ctx, "shell.RegisterShortcut", fn.Args{accelerator}, nil)
	return
}

// IsShortcutRegistered
func (m *ShellModule) IsShortcutRegistered(ctx context.Context, accelerator string) (ret bool, err error) {
	_, err = m.client.Call(ctx, "shell.IsShortcutRegistered", fn.Args{accelerator}, &ret)
	return
}

// UnregisterShortcut
func (m *ShellModule) UnregisterShortcut(ctx context.Context, accelerator string) (ret bool, err error) {
	_, err = m.client.Call(ctx, "shell.UnregisterShortcut", fn.Args{accelerator}, &ret)
	return
}

// UnregisterAllShortcuts
func (m *ShellModule) UnregisterAllShortcuts(ctx context.Context) (err error) {
	_, err = m.client.Call(ctx, "shell.UnregisterAllShortcuts", fn.Args{}, nil)
	return
}
