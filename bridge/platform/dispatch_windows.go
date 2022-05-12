package platform

func Dispatch(fn func()) {
	if !isRunning {
		fn()
		return
	}

	done := make(chan bool, 1)
	mainfunc <- func() {
		fn()
		done <- true
	}
	<-done
}
