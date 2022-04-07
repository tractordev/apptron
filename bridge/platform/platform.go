package platform

import "sync"

var (
	isRunning bool
	runReady  sync.WaitGroup
	once      sync.Once
)

func init() {
	runReady.Add(1)
}

func Start() {
	once.Do(func() {
		runReady.Done()
		isRunning = true
	})
}
