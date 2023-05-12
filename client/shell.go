package client

import (
	"context"
)

type ShellModule struct {
	client *Client
	OnShortcut func (event Event)
	ShowNotification func (ctx context.Context, n Notification) error
	ShowMessage func (ctx context.Context, msg MessageDialog) (bool, error)
	ShowFilePicker func (ctx context.Context, fd FileDialog) ([]string, error)
	ReadClipboard func (ctx context.Context) (string, error)
	WriteClipboard func (ctx context.Context, text string) (bool, error)
	RegisterShortcut func (ctx context.Context, accelerator string) error
	IsShortcutRegistered func (ctx context.Context, accelerator string) (bool, error)
	UnregisterShortcut func (ctx context.Context, accelerator string) (bool, error)
	UnregisterAllShortcuts func (ctx context.Context) error
}

type Notification struct {
	Title string
	Subtitle string
	Body string
}

type FileDialog struct {
	Title string
	Directory string
	Filename string
	Mode string
	Filters []string
}

type MessageDialog struct {
	Title string
	Body string
	Level string
	Buttons string
}

