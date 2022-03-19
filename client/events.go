package client

import (
	"fmt"
	"io"

	"github.com/progrium/qtalk-go/rpc"
)

type EventType int

const (
	EventNone EventType = iota
	EventClose
	EventDestroyed
	EventFocused
	EventBlurred
	EventResized
	EventMoved
	EventMenuItem
	EventShortcut
)

func (e EventType) String() string {
	return []string{"none", "close", "destroy", "focus", "blur", "resize", "move", "menu", "shortcut"}[e]
}

type Event struct {
	Type     EventType
	Name     string
	WindowID Handle
	Position Position
	Size     Size
	MenuID   uint16
	Shortcut string
}

func dispatchEvents(client *Client, resp *rpc.Response) {
	var e Event
	for {
		err := resp.Receive(&e)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err) // TODO: something else?
			}
			return
		}
		if client.OnEvent != nil {
			client.OnEvent(e)
		}
		switch e.Type {
		case EventMoved:
			w := client.Window.byID(e.WindowID)
			if w != nil && w.OnMoved != nil {
				w.OnMoved(e)
			}
		case EventResized:
			w := client.Window.byID(e.WindowID)
			if w != nil && w.OnResized != nil {
				w.OnResized(e)
			}
		case EventClose:
			w := client.Window.byID(e.WindowID)
			if w != nil && w.OnClose != nil {
				w.OnClose(e)
			}
		case EventDestroyed:
			w := client.Window.byID(e.WindowID)
			if w != nil && w.OnDestroyed != nil {
				w.OnDestroyed(e)
			}
		case EventFocused:
			w := client.Window.byID(e.WindowID)
			if w != nil && w.OnFocused != nil {
				w.OnFocused(e)
			}
		case EventBlurred:
			w := client.Window.byID(e.WindowID)
			if w != nil && w.OnBlurred != nil {
				w.OnBlurred(e)
			}
		case EventShortcut:
			if client.Shell.OnShortcut != nil {
				client.Shell.OnShortcut(e)
			}
		case EventMenuItem:
			if client.Menu.OnClick != nil {
				client.Menu.OnClick(e)
			}
		}
	}
}
