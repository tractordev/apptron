package client

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
