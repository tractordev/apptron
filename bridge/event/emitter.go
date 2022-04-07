package event

import "sync"

var listeners sync.Map

func Listen(key interface{}, cb func(e Event) error) {
	listeners.Store(key, cb)
}

func Unlisten(key interface{}) {
	listeners.Delete(key)
}

func Emit(event Event) {
	if event.Type <= 0 {
		return
	}
	listeners.Range(func(key, v interface{}) bool {
		cb := v.(func(Event) error)
		if err := cb(event); err != nil {
			listeners.Delete(key)
		}
		return true
	})
}
