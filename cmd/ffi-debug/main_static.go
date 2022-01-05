package main

// NOTE: There should be NO space between the comments and the `import "C"` line.
// The -ldl is necessary to fix the linker errors about `dlsym` that would otherwise appear.

/*
#cgo LDFLAGS: ./lib/libhostbridge.a -ldl
#include "../../lib/hostbridge.h"
*/
import "C"

func main() {
	C.hello(C.CString("static"))
	C.gomain()
}
