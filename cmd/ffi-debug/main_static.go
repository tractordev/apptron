package main

import "fmt"
import "github.com/progrium/hostbridge/bridge/window"

func main_loop(event window.Event) {
	if (event.Type > 0) {
		fmt.Println("[main_loop] event", event)

		if (event.Name == "close") {
			w := window.FindById(event.WindowId)
			if (w != nil) {
				w.Destroy()
			}

			all := window.All()
			fmt.Println("count of all windows", len(all))
			if (len(all) == 0) {
				fmt.Println("  quitting application...")
				window.Quit()
			}
		}
	}
}

func main() {
	options := window.Options{
		// NOTE(nick): resizing a transparent window on MacOS seems really slow?
		// Transparent: true,
		// Frameless: false,
		HTML: `
          <!doctype html>
          <html>
            <body style="font-family: -apple-system, BlinkMacSystemFont, avenir next, avenir, segoe ui, helvetica neue, helvetica, Ubuntu, roboto, noto, arial, sans-serif; background-color:rgba(87,87,87,0.8);"></body>
            <script>
              window.onload = function() {
                document.body.innerHTML = '<div style="padding: 30px">Transparency Test<br><br>${navigator.userAgent}</div>';
              };
            </script>
          </html>`,
	};

	w, _ := window.Create(options)

	fmt.Println("[main] window", w)

	w.SetTitle("Hello, Sailor!")
	fmt.Println("[main] window position", w.GetOuterPosition())

	w2, _ := window.Create(options)
	w2.SetTitle("YO!")

	w2.SetFullscreen(true)

	was_destroyed := w2.Destroy()
	fmt.Println("[main] was_destroyed", was_destroyed)

	window.Run(main_loop)

	// NOTE(nick): this doesn't appear to be called ever
	fmt.Println("[main] Goodbye.", w)
}
