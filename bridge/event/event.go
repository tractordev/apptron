package event

import (
	"tractor.dev/apptron/bridge/misc"
	"tractor.dev/apptron/bridge/resource"
)

var EventHandler func(event Event)

type Event struct {
	Type     EventType
	Window   resource.Handle
	Position misc.Position
	Size     misc.Size
	MenuItem int
	Shortcut string
}

type EventType int

const (
	None EventType = iota
	Close
	Created
	Destroyed
	Focused
	Blurred
	Resized
	Moved
	MenuItem
	Shortcut
)

func (e EventType) String() string {
	return []string{
		"",
		"close",
		"create",
		"destroy",
		"focus",
		"blur",
		"resize",
		"move",
		"menu",
		"shortcut",
	}[e]
}
