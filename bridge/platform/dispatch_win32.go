package platform

/*
import (
	"https://github.com/faiface/mainthread"
)
*/

func Dispatch(fn func()) {
	/*
	if !isRunning {
		fn()
		return
	}
	dispatch.Sync(dispatch.MainQueue(), fn)
	*/

	//fn()

	/*
	if !isRunning {
		fn()
		return
	}
	*/

	done := make(chan bool, 1)
	mainfunc <- func() {
		fn()
		done <- true
	}
	<-done
}
