package system

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/progrium/macdriver/cocoa"
)

func Displays() (displays []Display) {
	cs := cocoa.NSScreen_Screens()
	var screens []cocoa.NSScreen
	for i := uint64(0); i < cs.Count(); i++ {
		screens = append(screens, cocoa.NSScreen_FromRef(cs.ObjectAtIndex(i)))
	}
	for _, screen := range screens {
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

func Power() PowerInfo {
	result := PowerInfo{}

	// TODO: use native battery API with IOPowerSources.h from IOKit
	out, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		return result
	}

	// @Robustness: handle multiple batteries?
	lines := strings.Split(string(out), "\n")
	if len(lines) >= 1 {
		line := lines[1]

		endIndex := strings.Index(line, "%;")
		if endIndex >= 0 {
			startIndex := endIndex

			for startIndex >= 0 {
				if line[startIndex] == ' ' || line[startIndex] == '\t' {
					break
				}

				startIndex -= 1
			}

			if startIndex >= 0 {
				percentStr := line[startIndex+1 : endIndex]
				percent, err := strconv.Atoi(percentStr)
				if err == nil {
					result.BatteryPercent = float64(percent) / 100.0
				}
			}
		}

		result.IsOnBattery = !strings.Contains(line, " AC attached;") || strings.Contains(line, " discharging;")
		result.IsCharging = strings.Contains(line, " charging;")
	}

	return result
}
