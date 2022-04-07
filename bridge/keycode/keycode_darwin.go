package keycode

func init() {
	strToKeyCode["COMMAND"] = SuperLeft
	strToKeyCode["CMD"] = SuperLeft
	strToKeyCode["CMDLEFT"] = SuperLeft
	strToKeyCode["CMDRIGHT"] = SuperRight

	strToKeyCode["OPTION"] = AltLeft
	strToKeyCode["OPTIONLEFT"] = AltLeft
	strToKeyCode["OPTIONRIGHT"] = AltRight
}

func Scancode(k KeyCode) uint8 {
	return map[KeyCode]uint8{
		KeyA:            0x00,
		KeyS:            0x01,
		KeyD:            0x02,
		KeyF:            0x03,
		KeyH:            0x04,
		KeyG:            0x05,
		KeyZ:            0x06,
		KeyX:            0x07,
		KeyC:            0x08,
		KeyV:            0x09,
		KeyB:            0x0b,
		KeyQ:            0x0c,
		KeyW:            0x0d,
		KeyE:            0x0e,
		KeyR:            0x0f,
		KeyY:            0x10,
		KeyT:            0x11,
		Digit1:          0x12,
		Digit2:          0x13,
		Digit3:          0x14,
		Digit4:          0x15,
		Digit6:          0x16,
		Digit5:          0x17,
		Equal:           0x18,
		Digit9:          0x19,
		Digit7:          0x1a,
		Minus:           0x1b,
		Digit8:          0x1c,
		Digit0:          0x1d,
		BracketRight:    0x1e,
		KeyO:            0x1f,
		KeyU:            0x20,
		BracketLeft:     0x21,
		KeyI:            0x22,
		KeyP:            0x23,
		Enter:           0x24,
		KeyL:            0x25,
		KeyJ:            0x26,
		Quote:           0x27,
		KeyK:            0x28,
		Semicolon:       0x29,
		Backslash:       0x2a,
		Comma:           0x2b,
		Slash:           0x2c,
		KeyN:            0x2d,
		KeyM:            0x2e,
		Period:          0x2f,
		Tab:             0x30,
		Space:           0x31,
		Backquote:       0x32,
		Backspace:       0x33,
		Escape:          0x35,
		SuperRight:      0x36,
		SuperLeft:       0x37,
		ShiftLeft:       0x38,
		AltLeft:         0x3a,
		ControlLeft:     0x3b,
		ShiftRight:      0x3c,
		AltRight:        0x3d,
		ControlRight:    0x3e,
		F17:             0x40,
		NumpadDecimal:   0x41,
		NumpadMultiply:  0x43,
		NumpadAdd:       0x45,
		NumLock:         0x47,
		AudioVolumeUp:   0x49,
		AudioVolumeDown: 0x4a,
		NumpadDivide:    0x4b,
		NumpadEnter:     0x4c,
		NumpadSubtract:  0x4e,
		F18:             0x4f,
		F19:             0x50,
		NumpadEqual:     0x51,
		Numpad0:         0x52,
		Numpad1:         0x53,
		Numpad2:         0x54,
		Numpad3:         0x55,
		Numpad4:         0x56,
		Numpad5:         0x57,
		Numpad6:         0x58,
		Numpad7:         0x59,
		F20:             0x5a,
		Numpad8:         0x5b,
		Numpad9:         0x5c,
		IntlYen:         0x5d,
		F5:              0x60,
		F6:              0x61,
		F7:              0x62,
		F3:              0x63,
		F8:              0x64,
		F9:              0x65,
		F11:             0x67,
		F13:             0x69,
		F16:             0x6a,
		F14:             0x6b,
		F10:             0x6d,
		F12:             0x6f,
		F15:             0x71,
		Insert:          0x72,
		Home:            0x73,
		PageUp:          0x74,
		Delete:          0x75,
		F4:              0x76,
		End:             0x77,
		F2:              0x78,
		PageDown:        0x79,
		F1:              0x7a,
		ArrowLeft:       0x7b,
		ArrowRight:      0x7c,
		ArrowDown:       0x7d,
		ArrowUp:         0x7e,
	}[k]
}

func FromScancode(c uint8) KeyCode {
	return map[uint8]KeyCode{
		0x00: KeyA,
		0x01: KeyS,
		0x02: KeyD,
		0x03: KeyF,
		0x04: KeyH,
		0x05: KeyG,
		0x06: KeyZ,
		0x07: KeyX,
		0x08: KeyC,
		0x09: KeyV,
		//0x0a => World 1,
		0x0b: KeyB,
		0x0c: KeyQ,
		0x0d: KeyW,
		0x0e: KeyE,
		0x0f: KeyR,
		0x10: KeyY,
		0x11: KeyT,
		0x12: Digit1,
		0x13: Digit2,
		0x14: Digit3,
		0x15: Digit4,
		0x16: Digit6,
		0x17: Digit5,
		0x18: Equal,
		0x19: Digit9,
		0x1a: Digit7,
		0x1b: Minus,
		0x1c: Digit8,
		0x1d: Digit0,
		0x1e: BracketRight,
		0x1f: KeyO,
		0x20: KeyU,
		0x21: BracketLeft,
		0x22: KeyI,
		0x23: KeyP,
		0x24: Enter,
		0x25: KeyL,
		0x26: KeyJ,
		0x27: Quote,
		0x28: KeyK,
		0x29: Semicolon,
		0x2a: Backslash,
		0x2b: Comma,
		0x2c: Slash,
		0x2d: KeyN,
		0x2e: KeyM,
		0x2f: Period,
		0x30: Tab,
		0x31: Space,
		0x32: Backquote,
		0x33: Backspace,
		//0x34 => unknown,
		0x35: Escape,
		0x36: SuperRight,
		0x37: SuperLeft,
		0x38: ShiftLeft,
		0x39: CapsLock,
		0x3a: AltLeft,
		0x3b: ControlLeft,
		0x3c: ShiftRight,
		0x3d: AltRight,
		0x3e: ControlRight,
		0x3f: Fn,
		0x40: F17,
		0x41: NumpadDecimal,
		//0x42 -> unknown,
		0x43: NumpadMultiply,
		//0x44 => unknown,
		0x45: NumpadAdd,
		//0x46 => unknown,
		0x47: NumLock,
		//0x48: NumpadClear,

		// TODO: (Artur) for me, kVK_VolumeUp is 0x48
		// macOS 10.11
		// /System/Library/Frameworks/Carbon.framework/Versions/A/Frameworks/HIToolbox.framework/Versions/A/Headers/Events.h
		//0x49: AudioVolumeUp,
		0x49: AudioVolumeDown,
		0x4b: NumpadDivide,
		0x4c: NumpadEnter,
		//0x4d => unknown,
		0x4e: NumpadSubtract,
		0x4f: F18,
		0x50: F19,
		0x51: NumpadEqual,
		0x52: Numpad0,
		0x53: Numpad1,
		0x54: Numpad2,
		0x55: Numpad3,
		0x56: Numpad4,
		0x57: Numpad5,
		0x58: Numpad6,
		0x59: Numpad7,
		0x5a: F20,
		0x5b: Numpad8,
		0x5c: Numpad9,
		0x5d: IntlYen,
		//0x5e => JIS Ro,
		//0x5f => unknown,
		0x60: F5,
		0x61: F6,
		0x62: F7,
		0x63: F3,
		0x64: F8,
		0x65: F9,
		//0x66 => JIS Eisuu (macOS),
		0x67: F11,
		//0x68 => JIS Kanna (macOS),
		0x69: F13,
		0x6a: F16,
		0x6b: F14,
		//0x6c => unknown,
		0x6d: F10,
		//0x6e => unknown,
		0x6f: F12,
		//0x70 => unknown,
		0x71: F15,
		0x72: Insert,
		0x73: Home,
		0x74: PageUp,
		0x75: Delete,
		0x76: F4,
		0x77: End,
		0x78: F2,
		0x79: PageDown,
		0x7a: F1,
		0x7b: ArrowLeft,
		0x7c: ArrowRight,
		0x7d: ArrowDown,
		0x7e: ArrowUp,
		//0x7f =>  unknown,

		// 0xA is the caret (^) an macOS's German QERTZ layout. This key is at the same location as
		// backquote (`) on Windows' US layout.
		0xa: Backquote,
	}[c]
}
