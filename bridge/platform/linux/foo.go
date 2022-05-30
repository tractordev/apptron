package linux

// @Incomplete: see GLFW https://github.com/glfw/glfw/blob/master/src/x11_init.c#L1256

/*
type Library struct {
  handle unsafe.Pointer
}

func DLOpen(lib string) Library {
  libname := C.CString(lib)
  defer C.free(unsafe.Pointer(libname))

  var result Library

  result.handle = C.dlopen(libname, C.RTLD_LAZY | C.RTLD_LOCAL);
  if result.handle == nil {
    log.Println("[DLOpen] error opening library", lib)
    return result
  }

  return result
}

func DLSym(lib Library, procName string) func() {
  n := C.CString(procName)
  defer C.free(unsafe.Pointer(n))

  if lib.handle == nil {
    return nil
  }

  handle := C.dlsym(lib.handle, n)
  if handle == nil {
    return nil
  }

  log.Println("handle", handle)

  return *(*func())(unsafe.Pointer(handle))
}

func DLCall(fn func(), arg0 uintptr) uintptr {
  fn1 := *(*func(uintptr) uintptr)(unsafe.Pointer(&fn))
  log.Println("fn", fn)
  log.Println("fn1", fn1)
  return fn1(arg0)
}

func GoNuts() {
  x11 := DLOpen("libX11.so.6")

  log.Println(x11)

  XOpenDisplay := DLSym(x11, "XOpenDisplay")

  DLCall(XOpenDisplay, 0)
}
*/

