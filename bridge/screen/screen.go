package screen

/*
#include "../../lib/hostbridge.h"
*/
import "C"

import (
	"unsafe"

	"github.com/progrium/hostbridge/bridge/core"
)

var Module *module

func init() {
	Module = &module{}
}

type module struct{}

type Display struct {
	Name        string
	Size        Size
	Position    Position
	ScaleFactor float64
}

type Position struct {
	X float64
	Y float64
}

type Size struct {
	Width  float64
	Height float64
}

func Displays() []Display {
	return Module.Displays()
}

func (m module) Displays() []Display {
	ret := make(chan []Display)
	core.Dispatch(func() {
		array := C.screen_get_available_displays()

		n := int(array.count)
		result := make([]Display, n)

		items := (*[1 << 28]C.Display)(unsafe.Pointer(array.data))[:n:n]

		for i := 0; i < n; i++ {
			display := items[i]

			result[i] = Display{
				Name:        C.GoString(display.name),
				Size:        Size{Width: float64(display.size.width), Height: float64(display.size.height)},
				Position:    Position{X: float64(display.position.x), Y: float64(display.position.y)},
				ScaleFactor: float64(display.scale_factor),
			}
		}

		ret <- result
	})
	return <-ret
}
