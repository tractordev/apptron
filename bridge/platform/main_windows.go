package platform

import (
	"time"

	"tractor.dev/apptron/bridge/platform/win32"
)

var mainfunc = make(chan func())
var shouldQuit = false

func Main() {
	win32.OS_Init()

	runReady.Wait()

	// @Robustness: should we add some sort of framerate timing to this?
loop:
	for {
		win32.PollEvents()

		if shouldQuit {
			break loop
		}

		select {
		case fn := <-mainfunc:
			fn()
		default: // NOTE(nick): keep running at max speed!
		}

		// @Robustness: how accurate is time.Sleep? should we use win32.SleepMS instead?
		time.Sleep(1 * time.Millisecond)
	}

	win32.RemoveAllTrayMenus()

	win32.ExitProcess(0)
}

func Terminate(_dispatch bool) {
	shouldQuit = true
}
