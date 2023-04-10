package system

import (
	"tractor.dev/apptron/bridge/platform/win32"
)

func Displays() (displays []Display) {
	enumProc := func(monitor win32.HMONITOR, param1 win32.HDC, param2 *win32.RECT, param3 win32.LPARAM) uintptr {
		rect := param2
		info := win32.MONITORINFOEX{}
		devMode := win32.DEVMODE{}

		if win32.GetMonitorInfoW(monitor, &info) {
			if win32.EnumDisplaySettings(&info.DeviceName[0], win32.ENUM_CURRENT_SETTINGS, &devMode) {
				rect = &win32.RECT{
					Left:   devMode.DmPosition.X,
					Right:  devMode.DmPosition.X + win32.LONG(devMode.DmPelsWidth),
					Top:    devMode.DmPosition.Y,
					Bottom: devMode.DmPosition.Y + win32.LONG(devMode.DmPelsHeight),
				}
			}
		}

		/*
			cxLogicl := info.RcMonitor.Right - info.RcMonitor.Left
			cxPhysical := devMode.DmPelsWidth
			scaleFactor := float64(cxPhysical) / float64(cxLogicl)
		*/

		var dpiX win32.UINT
		var dpiY win32.UINT
		win32.GetDpiForMonitor(monitor, 0, &dpiX, &dpiY)

		scaleFactor := float64(dpiX) / float64(win32.USER_DEFAULT_SCREEN_DPI)

		displays = append(displays, Display{
			Name: info.GetDeviceName(),
			Size: Size{
				Width:  float64(rect.Right - rect.Left),
				Height: float64(rect.Bottom - rect.Top),
			},
			Position: Position{
				X: float64(rect.Top),
				Y: float64(rect.Left),
			},
			ScaleFactor: scaleFactor,
		})

		return uintptr(win32.TRUE)
	}

	win32.EnumDisplayMonitors(0, nil, enumProc, 0)

	return
}

func Power() PowerInfo {
	result := PowerInfo{}

	status := win32.SYSTEM_POWER_STATUS{}
	if win32.GetSystemPowerStatus(&status) {
		//
		// NOTE(nick): 255 indicates "unknown status" / failed to read battery information
		//
		// https://learn.microsoft.com/en-us/windows/win32/api/winbase/ns-winbase-system_power_status
		//

		if status.BatteryLifePercent != 255 {
			result.BatteryPercent = float64(status.BatteryLifePercent) / 100.0
		}

		if status.ACLineStatus != 255 {
			result.IsOnBattery = status.ACLineStatus != 1
		}

		if status.BatteryFlag != 255 {
			result.IsCharging = status.BatteryFlag&8 > 0
		}
	}

	return result
}
