package platform

func Main() {
	//app := cocoa.NSApp()
	//runReady.Wait()
	//app.Run()

loop:
	for {
		select {
		case f := <-mainfunc:
			f()
		case <-quit:
			break loop
		}
	}
}

var mainfunc = make(chan func())
var quit = make(chan bool)

func Terminate() {
	quit <- true
	/*
	Dispatch(func() {
		app := cocoa.NSApp()
		app.Terminate()
	})
	*/
}