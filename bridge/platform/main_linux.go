package platform

import (
	"time"

	"tractor.dev/apptron/bridge/platform/linux"
)

var mainfunc = make(chan func())
var shouldQuit = false

func init() {
	linux.OS_Init()
}

func Main() {
	runReady.Wait()

	// @Robustness: should we add some sort of framerate timing to this?
loop:
	for {
		linux.PollEvents()

		if shouldQuit {
			break loop
		}

		select {
		case fn := <-mainfunc:
			fn()
		default: // NOTE(nick): keep running at max speed!
		}

		// @Robustness: how accurate is time.Sleep? should we use linux.SleepMS instead?
		time.Sleep(1 * time.Millisecond)
	}

}

func Terminate(_dispatch bool) {
	shouldQuit = true
}
