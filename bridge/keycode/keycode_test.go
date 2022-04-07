package keycode

import "testing"

func TestKeyCodes(t *testing.T) {
	var code KeyCode
	code = KeyA
	if code.String() != "A" {
		t.Fatal("KeyA string was not 'A'")
	}
	if FromString("A") != KeyA {
		t.Fatal("'A' string not converted to KeyA")
	}
}
