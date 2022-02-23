package core

var EventHandler func(event Event)

type EventType int

const (
	EventNone EventType = iota
	EventClose
	EventDestroyed
	EventFocused
	EventResized
	EventMoved
	EventMenuItem
	EventShortcut
)

func (e EventType) String() string {
	return []string{"none", "close", "destroyed", "focused", "resized", "moved", "menu-item", "shortcut"}[e]
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
