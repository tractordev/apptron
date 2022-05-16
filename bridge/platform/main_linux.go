package platform

import (
	"time"

	"tractor.dev/apptron/bridge/platform/linux"
)

var mainfunc = make(chan func())
var quit = make(chan bool)

func init() {
	linux.OS_Init()
}

func Main() {
	runReady.Wait()

	// @Robustness: should we add some sort of framerate timing to this?
loop:
	for {
		linux.PollEvents()

		select {
		case fn := <-mainfunc:
			fn()
		case <-quit:
			break loop
		default: // NOTE(nick): keep running at max speed!
		}

		// @Robustness: how accurate is time.Sleep? should we use linux.SleepMS instead?
		time.Sleep(1 * time.Millisecond)
	}

}

func Terminate() {
	quit <- true
}
