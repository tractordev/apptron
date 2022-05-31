package system

import "tractor.dev/apptron/bridge/platform/linux"

func Displays() (displays []Display) {
  for _, monitor := range linux.Monitors() {
    rect := monitor.Geometry()

    displays = append(displays, Display{
      Name: monitor.Name(),
      Size: Size{
        Width:  float64(rect.Size.Width),
        Height: float64(rect.Size.Height),
      },
      Position: Position{
        X: float64(rect.Position.X),
        Y: float64(rect.Position.Y),
      },
      ScaleFactor: float64(monitor.ScaleFactor()),
    })
  }
  return
}
