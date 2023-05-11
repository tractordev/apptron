package client

import (
	"context"

	"github.com/progrium/qtalk-go/fn"
)

type SystemModule struct {
	client *Client
	Displays func () []Display
	Power func () PowerInfo
}

type Display struct {
	Name string
	Size Size
	Position Position
	ScaleFactor float64
}

type PowerInfo struct {
	IsOnBattery bool
	IsCharging bool
	BatteryPercent float64
}

