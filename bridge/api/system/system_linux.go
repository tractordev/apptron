package system

import (
  "os/exec"
  "strconv"
  "strings"

  "tractor.dev/apptron/bridge/platform/linux"
)

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

func Power() PowerInfo {
  result := PowerInfo{}

  out, err := exec.Command("upower", "-i", "/org/freedesktop/UPower/devices/battery_BAT0").Output()
  if err != nil {
    return result
  }

  str := string(out)

  charging := findKeyValue(str, "state:")
  result.IsCharging = charging == "charging"
  result.IsOnBattery = charging == "discharging"

  percentStr := findKeyValue(str, "percentage:")
  if len(percentStr) > 0 {
    if strings.HasSuffix(percentStr, "%") {
      percentStr = percentStr[:len(percentStr)-1]
    }
    percent, err := strconv.Atoi(percentStr)
    if err == nil {
      result.BatteryPercent = float64(percent) / 100.0
    }
  }

  return result
}

func findKeyValue(str string, key string) string {
  start := strings.Index(str, key)
  if start >= 0 {
    end := start + len(key) + 1

    for end < len(str) {
      if str[end] == '\n' {
        break
      }

      end += 1
    }

    if end < len(str) {
      return strings.TrimSpace(str[start+len(key) : end])
    }
  }

  return ""
}
