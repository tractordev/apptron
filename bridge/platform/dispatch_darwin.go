package platform

import "github.com/progrium/macdriver/dispatch"

func Dispatch(fn func()) {
	if !isRunning {
		fn()
		return
	}
	dispatch.Sync(dispatch.MainQueue(), fn)
}
