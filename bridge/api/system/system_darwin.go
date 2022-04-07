package system

import "github.com/progrium/macdriver/cocoa"

func Displays() (displays []Display) {
	for _, screen := range cocoa.NSScreen_Screens() {
		frame := screen.Frame()
		displays = append(displays, Display{
			Name: screen.LocalizedName().String(),
			Size: Size{
				Width:  frame.Size.Width,
				Height: frame.Size.Height,
			},
			Position: Position{
				X: frame.Origin.X,
				Y: frame.Origin.Y,
			},
			ScaleFactor: float64(screen.BackingScaleFactor()),
		})
	}
	return
}
