package keycode

func init() {
}

func Scancode(k KeyCode) uint8 {
  return map[KeyCode]uint8{
  }[k]
}

func FromScancode(c uint8) KeyCode {
  return map[uint8]KeyCode{
  }[c]
}
