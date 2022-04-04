package event

import (
	"tractor.dev/hostbridge/bridge/misc"
	"tractor.dev/hostbridge/bridge/resource"
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
		"destroy",
		"focus",
		"blur",
		"resize",
		"move",
		"menu",
		"shortcut",
	}[e]
}
