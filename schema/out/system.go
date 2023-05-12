package client

import (
	"context"
)

type SystemModule struct {
	client *Client
	Displays func (ctx context.Context) ([]Display, error)
	Power func (ctx context.Context) (PowerInfo, error)
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

