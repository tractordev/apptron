package apptron

import (
	"context"

	"tractor.dev/apptron/client"
)

// these aliases are experimental
// and may disappear

type Client = client.Client
type Window = client.Window
type Menu = client.Menu
type MenuItem = client.MenuItem
type FileDialog = client.FileDialog
type MessageDialog = client.MessageDialog
type Notification = client.Notification
type Display = client.Display
type AppOptions = client.AppOptions
type WindowOptions = client.WindowOptions
type Size = client.Size
type Position = client.Position
type Handle = client.Handle

func Run(ctx context.Context, opts AppOptions) (*Client, error) {
	client, err := client.Spawn()
	if err != nil {
		return nil, err
	}
	return client, client.App.Run(ctx, opts)
}
