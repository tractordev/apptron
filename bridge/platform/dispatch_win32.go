package platform

func Dispatch(fn func()) {
	/*
	if !isRunning {
		fn()
		return
	}
	dispatch.Sync(dispatch.MainQueue(), fn)
	*/

	fn()
}
