package main

import "fmt"
import "github.com/progrium/hostbridge/bridge/window"

func main_loop(event_type int) {
	if (event_type > 0) {
		fmt.Println("%d", event_type);
	}
}

func main() {
	w, _ := window.Create()

	fmt.Println("[main] window", w)

	w.SetTitle("Hello, Sailor!")
	fmt.Println("[main] window position", w.GetOuterPosition())

	w2, _ := window.Create()
	w2.SetTitle("YO!")

	window.Run(main_loop)
}
