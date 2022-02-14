package main

import "fmt"

import (
	"github.com/progrium/hostbridge/bridge/app"
	"github.com/progrium/hostbridge/bridge/menu"
	"github.com/progrium/hostbridge/bridge/window"
)

func tick(event app.Event) {
	if (event.Type > 0) {
		fmt.Println("[tick] event", event)

		if (event.Name == "close") {
			w := window.FindByID(event.WindowID)
			if (w != nil) {
				w.Destroy()
			}

			all := window.All()
			fmt.Println("count of all windows", len(all))
			if (len(all) == 0) {
				fmt.Println("  quitting application...")
				app.Quit()
			}
		}

		if (event.Name == "menu-item") {
			w := window.FindByID(event.WindowID)
			if (w != nil) {
				w.Destroy()
			}

			all := window.All()
			fmt.Println("count of all windows", len(all))
			if (len(all) == 0) {
				fmt.Println("  quitting application...")
				app.Quit()
			}
		}
	}
}

func main() {
	items := []menu.Item {
		{
			// NOTE(nick): when setting the window menu with wry, the first item title will always be the name of the executable on MacOS
			// so, this property is ignored:
			// @Robustness: maybe we want to make that more visible to the user somehow?
			Title: "this doesnt matter",
			Enabled: true,
			SubMenu: []menu.Item {
				{
					ID: 121,
					Title: "About",
					Enabled: true,
				},
				{
					ID: 122,
					Title: "Disabled",
					Enabled: false,
				},
				{
					ID: 99,
					Title: "Quit",
					Enabled: true,
					//Accelerator: "CommandOrControl+Q",
				},
			},
		},
		{
			ID: 23,
			Title: "hello world",
			Enabled: true,
			SubMenu: []menu.Item {
				{
					ID: 777,
					Title: "About2",
					Enabled: true,
				},
			},
		},
	}

	m := menu.New(items)
	app.SetMenu(m)

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
			</html>
		`,
	};

	w1, _ := window.Create(options)

	fmt.Println("[main] window", w1)

	if (w1 == nil) {
		return
	}

	w1.SetTitle("Hello, Sailor!")
	fmt.Println("[main] window position", w1.GetOuterPosition())

	w2, _ := window.Create(options)
	w2.SetTitle("YO!")

	w2.SetFullscreen(true)

	wasDestroyed := w2.Destroy()
	fmt.Println("[main] wasDestroyed", wasDestroyed)

	app.Run(tick)

	// NOTE(nick): this doesn't appear to be called ever
	fmt.Println("[main] Goodbye.", w1)
}
