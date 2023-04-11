package client

import (
	"context"

	"github.com/progrium/qtalk-go/fn"
)

type Display struct {
	Name        string
	Size        Size
	Position    Position
	ScaleFactor float64
}

type PowerInfo struct {
	IsOnBattery    bool
	IsCharging     bool
	BatteryPercent float64
}

type SystemModule struct {
	client *Client
}

// Displays
func (m *SystemModule) Displays(ctx context.Context) (ret []Display, err error) {
	_, err = m.client.Call(ctx, "system.Displays", fn.Args{}, &ret)
	return
}

// Power
func (m *SystemModule) Power(ctx context.Context) (ret PowerInfo, err error) {
	_, err = m.client.Call(ctx, "system.Power", fn.Args{}, &ret)
	return
}
