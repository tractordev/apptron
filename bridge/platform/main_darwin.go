package platform

import (
	"github.com/progrium/macdriver/cocoa"
)

func Main() {
	app := cocoa.NSApp()
	runReady.Wait()
	app.Run()
}

func Terminate(dispatch bool) {
	fn := func() {
		app := cocoa.NSApp()
		app.Terminate()
	}
	if dispatch {
		Dispatch(fn)
	} else {
		fn()
	}
}
