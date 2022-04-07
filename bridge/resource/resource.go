package resource

import (
	"errors"
	"sync"

	"github.com/rs/xid"
)

var (
	ErrBadHandle = errors.New("bad handle")

	resources sync.Map
	released  sync.Map
)

type Handle string

func NewHandle() Handle {
	return Handle(xid.New().String())
}

func Retain(h Handle, v interface{}) {
	resources.Store(h, v)
}

func Release(h Handle) {
	resources.Delete(h)
	released.Store(h, struct{}{})
}

func IsReleased(h Handle) (found bool) {
	_, found = released.Load(h)
	return
}

func Lookup(h Handle) (interface{}, error) {
	v, found := resources.Load(h)
	if !found {
		return nil, ErrBadHandle
	}
	return v, nil
}

func Range(fn func(interface{}) bool) {
	resources.Range(func(key, value interface{}) bool {
		return fn(value)
	})
}
