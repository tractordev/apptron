package system

import "tractor.dev/hostbridge/bridge/misc"

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

type Position = misc.Position

type Size = misc.Size

func (m module) Displays() []Display {
	return Displays()
}
