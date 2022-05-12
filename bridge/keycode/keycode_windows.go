package keycode

import "golang.design/x/hotkey"

func init() {
	strToKeyCode["WINDOWS"] = SuperLeft
	strToKeyCode["WIN"] = SuperLeft
	strToKeyCode["WINLEFT"] = SuperLeft
	strToKeyCode["WINRIGHT"] = SuperRight
}

func HotkeyModifier(code KeyCode) hotkey.Modifier {
	return map[KeyCode]hotkey.Modifier{
		AltLeft:      hotkey.ModAlt,
		AltRight:     hotkey.ModAlt,
		ControlLeft:  hotkey.ModCtrl,
		ControlRight: hotkey.ModCtrl,
		ShiftLeft:    hotkey.ModShift,
		ShiftRight:   hotkey.ModShift,
		SuperLeft:    hotkey.ModWin,
		SuperRight:   hotkey.ModWin,
	}[code]
}

func Scancode(k KeyCode) uint8 {
	return map[KeyCode]uint8{}[k]
}

func FromScancode(c uint8) KeyCode {
	return map[uint8]KeyCode{}[c]
}
