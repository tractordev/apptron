// This file contains a substantial portion of the UI Events Specification by the W3C. In
// particular, the variant names within `KeyCode` and their documentation are modified
// versions of contents of the aforementioned specification.
//
// The original documents are:
//
// ### For `KeyCode`
// UI Events KeyboardEvent code Values
// https://www.w3.org/TR/2017/CR-uievents-code-20170601/
// Copyright © 2017 W3C® (MIT, ERCIM, Keio, Beihang).
//
// These documents were used under the terms of the following license. This W3C license as well as
// the W3C short notice apply to the `KeyCode` enums and their variants and the
// documentation attached to their variants.

// --------- BEGGINING OF W3C LICENSE --------------------------------------------------------------
//
// License
//
// By obtaining and/or copying this work, you (the licensee) agree that you have read, understood,
// and will comply with the following terms and conditions.
//
// Permission to copy, modify, and distribute this work, with or without modification, for any
// purpose and without fee or royalty is hereby granted, provided that you include the following on
// ALL copies of the work or portions thereof, including modifications:
//
// - The full text of this NOTICE in a location viewable to users of the redistributed or derivative
//   work.
// - Any pre-existing intellectual property disclaimers, notices, or terms and conditions. If none
//   exist, the W3C Software and Document Short Notice should be included.
// - Notice of any changes or modifications, through a copyright statement on the new code or
//   document such as "This software or document includes material copied from or derived from
//   [title and URI of the W3C document]. Copyright © [YEAR] W3C® (MIT, ERCIM, Keio, Beihang)."
//
// Disclaimers
//
// THIS WORK IS PROVIDED "AS IS," AND COPYRIGHT HOLDERS MAKE NO REPRESENTATIONS OR WARRANTIES,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO, WARRANTIES OF MERCHANTABILITY OR FITNESS FOR
// ANY PARTICULAR PURPOSE OR THAT THE USE OF THE SOFTWARE OR DOCUMENT WILL NOT INFRINGE ANY THIRD
// PARTY PATENTS, COPYRIGHTS, TRADEMARKS OR OTHER RIGHTS.
//
// COPYRIGHT HOLDERS WILL NOT BE LIABLE FOR ANY DIRECT, INDIRECT, SPECIAL OR CONSEQUENTIAL DAMAGES
// ARISING OUT OF ANY USE OF THE SOFTWARE OR DOCUMENT.
//
// The name and trademarks of copyright holders may NOT be used in advertising or publicity
// pertaining to the work without specific, written prior permission. Title to copyright in this
// work will at all times remain with copyright holders.
//
// --------- END OF W3C LICENSE --------------------------------------------------------------------

// --------- BEGGINING OF W3C SHORT NOTICE ---------------------------------------------------------
//
// tao: https://github.com/tauri-apps/tao
//
// Copyright © 2021 World Wide Web Consortium, (Massachusetts Institute of Technology, European
// Research Consortium for Informatics and Mathematics, Keio University, Beihang). All Rights
// Reserved. This work is distributed under the W3C® Software License [1] in the hope that it will
// be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// [1] http://www.w3.org/Consortium/Legal/copyright-software
//
// --------- END OF W3C SHORT NOTICE ---------------------------------------------------------------
package keycode

type KeyCode uint

const (
	Unknown KeyCode = iota
	/// <kbd>`</kbd> on a US keyboard. This is also called a backtick or grave.
	/// This is the <kbd>半角</kbd>/<kbd>全角</kbd>/<kbd>漢字</kbd>
	/// (hankaku/zenkaku/kanji) key on Japanese keyboards
	Backquote
	/// Used for both the US <kbd>\\</kbd> (on the 101-key layout) and also for the key
	/// located between the <kbd>"</kbd> and <kbd>Enter</kbd> keys on row C of the 102-,
	/// 104- and 106-key layouts.
	/// Labeled <kbd>#</kbd> on a UK (102) keyboard.
	Backslash
	/// <kbd>[</kbd> on a US keyboard.
	BracketLeft
	/// <kbd>]</kbd> on a US keyboard.
	BracketRight
	/// <kbd></kbd> on a US keyboard.
	Comma
	/// <kbd>0</kbd> on a US keyboard.
	Digit0
	/// <kbd>1</kbd> on a US keyboard.
	Digit1
	/// <kbd>2</kbd> on a US keyboard.
	Digit2
	/// <kbd>3</kbd> on a US keyboard.
	Digit3
	/// <kbd>4</kbd> on a US keyboard.
	Digit4
	/// <kbd>5</kbd> on a US keyboard.
	Digit5
	/// <kbd>6</kbd> on a US keyboard.
	Digit6
	/// <kbd>7</kbd> on a US keyboard.
	Digit7
	/// <kbd>8</kbd> on a US keyboard.
	Digit8
	/// <kbd>9</kbd> on a US keyboard.
	Digit9
	/// <kbd>=</kbd> on a US keyboard.
	Equal
	/// Located between the left <kbd>Shift</kbd> and <kbd>Z</kbd> keys.
	/// Labeled <kbd>\\</kbd> on a UK keyboard.
	IntlBackslash
	/// Located between the <kbd>/</kbd> and right <kbd>Shift</kbd> keys.
	/// Labeled <kbd>\\</kbd> (ro) on a Japanese keyboard.
	IntlRo
	/// Located between the <kbd>=</kbd> and <kbd>Backspace</kbd> keys.
	/// Labeled <kbd>¥</kbd> (yen) on a Japanese keyboard. <kbd>\\</kbd> on a
	/// Russian keyboard.
	IntlYen
	/// <kbd>a</kbd> on a US keyboard.
	/// Labeled <kbd>q</kbd> on an AZERTY (e.g. French) keyboard.
	KeyA
	/// <kbd>b</kbd> on a US keyboard.
	KeyB
	/// <kbd>c</kbd> on a US keyboard.
	KeyC
	/// <kbd>d</kbd> on a US keyboard.
	KeyD
	/// <kbd>e</kbd> on a US keyboard.
	KeyE
	/// <kbd>f</kbd> on a US keyboard.
	KeyF
	/// <kbd>g</kbd> on a US keyboard.
	KeyG
	/// <kbd>h</kbd> on a US keyboard.
	KeyH
	/// <kbd>i</kbd> on a US keyboard.
	KeyI
	/// <kbd>j</kbd> on a US keyboard.
	KeyJ
	/// <kbd>k</kbd> on a US keyboard.
	KeyK
	/// <kbd>l</kbd> on a US keyboard.
	KeyL
	/// <kbd>m</kbd> on a US keyboard.
	KeyM
	/// <kbd>n</kbd> on a US keyboard.
	KeyN
	/// <kbd>o</kbd> on a US keyboard.
	KeyO
	/// <kbd>p</kbd> on a US keyboard.
	KeyP
	/// <kbd>q</kbd> on a US keyboard.
	/// Labeled <kbd>a</kbd> on an AZERTY (e.g. French) keyboard.
	KeyQ
	/// <kbd>r</kbd> on a US keyboard.
	KeyR
	/// <kbd>s</kbd> on a US keyboard.
	KeyS
	/// <kbd>t</kbd> on a US keyboard.
	KeyT
	/// <kbd>u</kbd> on a US keyboard.
	KeyU
	/// <kbd>v</kbd> on a US keyboard.
	KeyV
	/// <kbd>w</kbd> on a US keyboard.
	/// Labeled <kbd>z</kbd> on an AZERTY (e.g. French) keyboard.
	KeyW
	/// <kbd>x</kbd> on a US keyboard.
	KeyX
	/// <kbd>y</kbd> on a US keyboard.
	/// Labeled <kbd>z</kbd> on a QWERTZ (e.g. German) keyboard.
	KeyY
	/// <kbd>z</kbd> on a US keyboard.
	/// Labeled <kbd>w</kbd> on an AZERTY (e.g. French) keyboard and <kbd>y</kbd> on a
	/// QWERTZ (e.g. German) keyboard.
	KeyZ
	/// <kbd>-</kbd> on a US keyboard.
	Minus
	/// <kbd>.</kbd> on a US keyboard.
	Period
	/// <kbd>'</kbd> on a US keyboard.
	Quote
	/// <kbd>;</kbd> on a US keyboard.
	Semicolon
	/// <kbd>/</kbd> on a US keyboard.
	Slash
	/// <kbd>Alt</kbd> <kbd>Option</kbd> or <kbd>⌥</kbd>.
	AltLeft
	/// <kbd>Alt</kbd> <kbd>Option</kbd> or <kbd>⌥</kbd>.
	/// This is labeled <kbd>AltGr</kbd> on many keyboard layouts.
	AltRight
	/// <kbd>Backspace</kbd> or <kbd>⌫</kbd>.
	/// Labeled <kbd>Delete</kbd> on Apple keyboards.
	Backspace
	/// <kbd>CapsLock</kbd> or <kbd>⇪</kbd>
	CapsLock
	/// The application context menu key which is typically found between the right
	/// <kbd>Super</kbd> key and the right <kbd>Control</kbd> key.
	ContextMenu
	/// <kbd>Control</kbd> or <kbd>⌃</kbd>
	ControlLeft
	/// <kbd>Control</kbd> or <kbd>⌃</kbd>
	ControlRight
	/// <kbd>Enter</kbd> or <kbd>↵</kbd>. Labeled <kbd>Return</kbd> on Apple keyboards.
	Enter
	/// The Windows <kbd>⌘</kbd> <kbd>Command</kbd> or other OS symbol key.
	SuperLeft
	/// The Windows <kbd>⌘</kbd> <kbd>Command</kbd> or other OS symbol key.
	SuperRight
	/// <kbd>Shift</kbd> or <kbd>⇧</kbd>
	ShiftLeft
	/// <kbd>Shift</kbd> or <kbd>⇧</kbd>
	ShiftRight
	/// <kbd> </kbd> (space)
	Space
	/// <kbd>Tab</kbd> or <kbd>⇥</kbd>
	Tab
	/// Japanese: <kbd>変</kbd> (henkan)
	Convert
	/// Japanese: <kbd>カタカナ</kbd>/<kbd>ひらがな</kbd>/<kbd>ローマ字</kbd> (katakana/hiragana/romaji)
	KanaMode
	/// Korean: HangulMode <kbd>한/영</kbd> (han/yeong)
	///
	/// Japanese (Mac keyboard): <kbd>か</kbd> (kana)
	Lang1
	/// Korean: Hanja <kbd>한</kbd> (hanja)
	///
	/// Japanese (Mac keyboard): <kbd>英</kbd> (eisu)
	Lang2
	/// Japanese (word-processing keyboard): Katakana
	Lang3
	/// Japanese (word-processing keyboard): Hiragana
	Lang4
	/// Japanese (word-processing keyboard): Zenkaku/Hankaku
	Lang5
	/// Japanese: <kbd>無変換</kbd> (muhenkan)
	NonConvert
	/// <kbd>⌦</kbd>. The forward delete key.
	/// Note that on Apple keyboards the key labelled <kbd>Delete</kbd> on the main part of
	/// the keyboard is encoded as [`Backspace`].
	///
	/// [`Backspace`]: Self::Backspace
	Delete
	/// <kbd>Page Down</kbd> <kbd>End</kbd> or <kbd>↘</kbd>
	End
	/// <kbd>Help</kbd>. Not present on standard PC keyboards.
	Help
	/// <kbd>Home</kbd> or <kbd>↖</kbd>
	Home
	/// <kbd>Insert</kbd> or <kbd>Ins</kbd>. Not present on Apple keyboards.
	Insert
	/// <kbd>Page Down</kbd> <kbd>PgDn</kbd> or <kbd>⇟</kbd>
	PageDown
	/// <kbd>Page Up</kbd> <kbd>PgUp</kbd> or <kbd>⇞</kbd>
	PageUp
	/// <kbd>↓</kbd>
	ArrowDown
	/// <kbd>←</kbd>
	ArrowLeft
	/// <kbd>→</kbd>
	ArrowRight
	/// <kbd>↑</kbd>
	ArrowUp
	/// On the Mac this is used for the numpad <kbd>Clear</kbd> key.
	NumLock
	/// <kbd>0 Ins</kbd> on a keyboard. <kbd>0</kbd> on a phone or remote control
	Numpad0
	/// <kbd>1 End</kbd> on a keyboard. <kbd>1</kbd> or <kbd>1 QZ</kbd> on a phone or remote control
	Numpad1
	/// <kbd>2 ↓</kbd> on a keyboard. <kbd>2 ABC</kbd> on a phone or remote control
	Numpad2
	/// <kbd>3 PgDn</kbd> on a keyboard. <kbd>3 DEF</kbd> on a phone or remote control
	Numpad3
	/// <kbd>4 ←</kbd> on a keyboard. <kbd>4 GHI</kbd> on a phone or remote control
	Numpad4
	/// <kbd>5</kbd> on a keyboard. <kbd>5 JKL</kbd> on a phone or remote control
	Numpad5
	/// <kbd>6 →</kbd> on a keyboard. <kbd>6 MNO</kbd> on a phone or remote control
	Numpad6
	/// <kbd>7 Home</kbd> on a keyboard. <kbd>7 PQRS</kbd> or <kbd>7 PRS</kbd> on a phone
	/// or remote control
	Numpad7
	/// <kbd>8 ↑</kbd> on a keyboard. <kbd>8 TUV</kbd> on a phone or remote control
	Numpad8
	/// <kbd>9 PgUp</kbd> on a keyboard. <kbd>9 WXYZ</kbd> or <kbd>9 WXY</kbd> on a phone
	/// or remote control
	Numpad9
	/// <kbd>+</kbd>
	NumpadAdd
	/// Found on the Microsoft Natural Keyboard.
	NumpadBackspace
	/// <kbd>C</kbd> or <kbd>A</kbd> (All Clear). Also for use with numpads that have a
	/// <kbd>Clear</kbd> key that is separate from the <kbd>NumLock</kbd> key. On the Mac the
	/// numpad <kbd>Clear</kbd> key is encoded as [`NumLock`].
	///
	/// [`NumLock`]: Self::NumLock
	NumpadClear
	/// <kbd>C</kbd> (Clear Entry)
	NumpadClearEntry
	/// <kbd></kbd> (thousands separator). For locales where the thousands separator
	/// is a "." (e.g. Brazil) this key may generate a <kbd>.</kbd>.
	NumpadComma
	/// <kbd>. Del</kbd>. For locales where the decimal separator is "" (e.g.
	/// Brazil) this key may generate a <kbd></kbd>.
	NumpadDecimal
	/// <kbd>/</kbd>
	NumpadDivide
	NumpadEnter
	/// <kbd>=</kbd>
	NumpadEqual
	/// <kbd>#</kbd> on a phone or remote control device. This key is typically found
	/// below the <kbd>9</kbd> key and to the right of the <kbd>0</kbd> key.
	NumpadHash
	/// <kbd>M</kbd> Add current entry to the value stored in memory.
	NumpadMemoryAdd
	/// <kbd>M</kbd> Clear the value stored in memory.
	NumpadMemoryClear
	/// <kbd>M</kbd> Replace the current entry with the value stored in memory.
	NumpadMemoryRecall
	/// <kbd>M</kbd> Replace the value stored in memory with the current entry.
	NumpadMemoryStore
	/// <kbd>M</kbd> Subtract current entry from the value stored in memory.
	NumpadMemorySubtract
	/// <kbd>*</kbd> on a keyboard. For use with numpads that provide mathematical
	/// operations (<kbd>+</kbd> <kbd>-</kbd> <kbd>*</kbd> and <kbd>/</kbd>).
	///
	/// Use `NumpadStar` for the <kbd>*</kbd> key on phones and remote controls.
	NumpadMultiply
	/// <kbd>(</kbd> Found on the Microsoft Natural Keyboard.
	NumpadParenLeft
	/// <kbd>)</kbd> Found on the Microsoft Natural Keyboard.
	NumpadParenRight
	/// <kbd>*</kbd> on a phone or remote control device.
	///
	/// This key is typically found below the <kbd>7</kbd> key and to the left of
	/// the <kbd>0</kbd> key.
	///
	/// Use <kbd>"NumpadMultiply"</kbd> for the <kbd>*</kbd> key on
	/// numeric keypads.
	NumpadStar
	/// <kbd>-</kbd>
	NumpadSubtract
	/// <kbd>Esc</kbd> or <kbd>⎋</kbd>
	Escape
	/// <kbd>Fn</kbd> This is typically a hardware key that does not generate a separate code.
	Fn
	/// <kbd>FLock</kbd> or <kbd>FnLock</kbd>. Function Lock key. Found on the Microsoft
	/// Natural Keyboard.
	FnLock
	/// <kbd>PrtScr SysRq</kbd> or <kbd>Print Screen</kbd>
	PrintScreen
	/// <kbd>Scroll Lock</kbd>
	ScrollLock
	/// <kbd>Pause Break</kbd>
	Pause
	/// Some laptops place this key to the left of the <kbd>↑</kbd> key.
	///
	/// This also the "back" button (triangle) on Android.
	BrowserBack
	BrowserFavorites
	/// Some laptops place this key to the right of the <kbd>↑</kbd> key.
	BrowserForward
	/// The "home" button on Android.
	BrowserHome
	BrowserRefresh
	BrowserSearch
	BrowserStop
	/// <kbd>Eject</kbd> or <kbd>⏏</kbd>. This key is placed in the function section on some Apple
	/// keyboards.
	Eject
	/// Sometimes labelled <kbd>My Computer</kbd> on the keyboard
	LaunchApp1
	/// Sometimes labelled <kbd>Calculator</kbd> on the keyboard
	LaunchApp2
	LaunchMail
	MediaPlayPause
	MediaSelect
	MediaStop
	MediaTrackNext
	MediaTrackPrevious
	/// This key is placed in the function section on some Apple keyboards replacing the
	/// <kbd>Eject</kbd> key.
	Power
	Sleep
	AudioVolumeDown
	AudioVolumeMute
	AudioVolumeUp
	WakeUp
	Hyper
	Turbo
	Abort
	Resume
	Suspend
	/// Found on Sun’s USB keyboard.
	Again
	/// Found on Sun’s USB keyboard.
	Copy
	/// Found on Sun’s USB keyboard.
	Cut
	/// Found on Sun’s USB keyboard.
	Find
	/// Found on Sun’s USB keyboard.
	Open
	/// Found on Sun’s USB keyboard.
	Paste
	/// Found on Sun’s USB keyboard.
	Props
	/// Found on Sun’s USB keyboard.
	Select
	/// Found on Sun’s USB keyboard.
	Undo
	/// Use for dedicated <kbd>ひらがな</kbd> key found on some Japanese word processing keyboards.
	Hiragana
	/// Use for dedicated <kbd>カタカナ</kbd> key found on some Japanese word processing keyboards.
	Katakana
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F1
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F2
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F3
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F4
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F5
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F6
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F7
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F8
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F9
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F10
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F11
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F12
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F13
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F14
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F15
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F16
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F17
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F18
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F19
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F20
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F21
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F22
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F23
	/// General-purpose function key.
	/// Usually found at the top of the keyboard.
	F24
	/// General-purpose function key.
	F25
	/// General-purpose function key.
	F26
	/// General-purpose function key.
	F27
	/// General-purpose function key.
	F28
	/// General-purpose function key.
	F29
	/// General-purpose function key.
	F30
	/// General-purpose function key.
	F31
	/// General-purpose function key.
	F32
	/// General-purpose function key.
	F33
	/// General-purpose function key.
	F34
	/// General-purpose function key.
	F35
)

var keyCodeToStr map[KeyCode]string
var strToKeyCode = map[string]KeyCode{
	"ALT":        AltLeft,
	"ALTLEFT":    AltLeft,
	"ALTRIGHT":   AltRight,
	"SHIFT":      ShiftLeft,
	"SHIFTLEFT":  ShiftLeft,
	"SHIFTRIGHT": ShiftRight,
	"CONTROL":    ControlLeft,
	"CTRL":       ControlLeft,
	"CTRLLEFT":   ControlLeft,
	"CTRLRIGHT":  ControlRight,
	"SUPER":      SuperLeft,
	"SUPERLEFT":  SuperLeft,
	"SUPERRIGHT": SuperRight,

	"`":            Backquote,
	"BACKQUOTE":    Backquote,
	"BACKSLASH":    Backslash,
	"[":            BracketLeft,
	"BRACKETLEFT":  BracketLeft,
	"]":            BracketRight,
	"BRACKETRIGHT": BracketRight,
	",":            Comma,
	"COMMA":        Comma,
	"0":            Digit0,
	"1":            Digit1,
	"2":            Digit2,
	"3":            Digit3,
	"4":            Digit4,
	"5":            Digit5,
	"6":            Digit6,
	"7":            Digit7,
	"8":            Digit8,
	"9":            Digit9,
	"NUM0":         Numpad0,
	"NUMPAD0":      Numpad0,
	"NUM1":         Numpad1,
	"NUMPAD1":      Numpad1,
	"NUM2":         Numpad2,
	"NUMPAD2":      Numpad2,
	"NUM3":         Numpad3,
	"NUMPAD3":      Numpad3,
	"NUM4":         Numpad4,
	"NUMPAD4":      Numpad4,
	"NUM5":         Numpad5,
	"NUMPAD5":      Numpad5,
	"NUM6":         Numpad6,
	"NUMPAD6":      Numpad6,
	"NUM7":         Numpad7,
	"NUMPAD7":      Numpad7,
	"NUM8":         Numpad8,
	"NUMPAD8":      Numpad8,
	"NUM9":         Numpad9,
	"NUMPAD9":      Numpad9,
	"=":            Equal,
	"-":            Minus,
	".":            Period,
	"PERIOD":       Period,
	"'":            Quote,
	"QUOTE":        Quote,
	"\\":           IntlBackslash,
	"A":            KeyA,
	"B":            KeyB,
	"C":            KeyC,
	"D":            KeyD,
	"E":            KeyE,
	"F":            KeyF,
	"G":            KeyG,
	"H":            KeyH,
	"I":            KeyI,
	"J":            KeyJ,
	"K":            KeyK,
	"L":            KeyL,
	"M":            KeyM,
	"N":            KeyN,
	"O":            KeyO,
	"P":            KeyP,
	"Q":            KeyQ,
	"R":            KeyR,
	"S":            KeyS,
	"T":            KeyT,
	"U":            KeyU,
	"V":            KeyV,
	"W":            KeyW,
	"X":            KeyX,
	"Y":            KeyY,
	"Z":            KeyZ,

	";":           Semicolon,
	"SEMICOLON":   Semicolon,
	"/":           Slash,
	"SLASH":       Slash,
	"BACKSPACE":   Backspace,
	"CAPSLOCK":    CapsLock,
	"CONTEXTMENU": ContextMenu,
	"ENTER":       Enter,
	"SPACE":       Space,
	"TAB":         Tab,
	"CONVERT":     Convert,

	"DELETE":   Delete,
	"END":      End,
	"HELP":     Help,
	"HOME":     Home,
	"PAGEDOWN": PageDown,
	"PAGEUP":   PageUp,

	"DOWN":  ArrowDown,
	"UP":    ArrowUp,
	"LEFT":  ArrowLeft,
	"RIGHT": ArrowRight,

	"NUMLOCK":         NumLock,
	"NUMADD":          NumpadAdd,
	"NUMPADADD":       NumpadAdd,
	"NUMBACKSPACE":    NumpadBackspace,
	"NUMPADBACKSPACE": NumpadBackspace,
	"NUMCLEAR":        NumpadClear,
	"NUMPADCLEAR":     NumpadClear,
	"NUMCOMMA":        NumpadComma,
	"NUMPADCOMMA":     NumpadComma,
	"NUMDIVIDE":       NumpadDivide,
	"NUMPADDIVIDE":    NumpadDivide,
	"NUMSUBSTRACT":    NumpadSubtract,
	"NUMPADSUBSTRACT": NumpadSubtract,
	"NUMENTER":        NumpadEnter,
	"NUMPADENTER":     NumpadEnter,

	"ESC":         Escape,
	"ESCAPE":      Escape,
	"FN":          Fn,
	"FNLOCK":      FnLock,
	"PRINTSCREEN": PrintScreen,
	"SCROLLLOCK":  ScrollLock,

	"PAUSE": Pause,

	"VOLUMEMUTE":         AudioVolumeMute,
	"VOLUMEDOWN":         AudioVolumeDown,
	"VOLUMEUP":           AudioVolumeUp,
	"MEDIANEXTTRACK":     MediaTrackNext,
	"MEDIAPREVIOUSTRACK": MediaTrackPrevious,
	"MEDIAPLAYPAUSE":     MediaPlayPause,
	"LAUNCHMAIL":         LaunchMail,

	"SUSPEND": Suspend,
	"F1":      F1,
	"F2":      F2,
	"F3":      F3,
	"F4":      F4,
	"F5":      F5,
	"F6":      F6,
	"F7":      F7,
	"F8":      F8,
	"F9":      F9,
	"F10":     F10,
	"F11":     F11,
	"F12":     F12,
	"F13":     F13,
	"F14":     F14,
	"F15":     F15,
	"F16":     F16,
	"F17":     F17,
	"F18":     F18,
	"F19":     F19,
	"F20":     F20,
	"F21":     F21,
	"F22":     F22,
	"F23":     F23,
	"F24":     F24,
	"F25":     F25,
	"F26":     F26,
	"F27":     F27,
	"F28":     F28,
	"F29":     F29,
	"F30":     F30,
	"F31":     F31,
	"F32":     F32,
	"F33":     F33,
	"F34":     F34,
	"F35":     F35,
}

func init() {
	keyCodeToStr = make(map[KeyCode]string)
	for k, v := range strToKeyCode {
		keyCodeToStr[v] = k
	}
}

func FromString(s string) KeyCode {
	return strToKeyCode[s]
}

func (c KeyCode) String() string {
	return keyCodeToStr[c]
}

func IsModifier(c KeyCode) bool {
	switch c {
	case ShiftLeft, ShiftRight, ControlLeft, ControlRight, AltLeft, AltRight, SuperLeft, SuperRight:
		return true
	default:
		return false
	}
}
