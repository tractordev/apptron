//go:build pkg

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"

	"github.com/progrium/macdriver/core"
	"tractor.dev/hostbridge/bridge/api/app"
	"tractor.dev/hostbridge/bridge/api/menu"
	"tractor.dev/hostbridge/bridge/api/system"
	"tractor.dev/hostbridge/bridge/api/window"
	"tractor.dev/hostbridge/bridge/platform"
)

func init() {
	runtime.LockOSThread()
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	go run()
	platform.Main()
}

func run() {
	items := []menu.Item{
		{
			Title: "App",
			SubMenu: []menu.Item{
				{
					Title: "One",
				},
				{
					Title: "Two",
				},
			},
		},
		{
			Title: "File",
			SubMenu: []menu.Item{
				{
					Title: "ONE",
				},
				{
					Title: "TWO",
				},
			},
		},
		{
			Title: "Foo",
		},
	}
	m := menu.New(items)
	app.SetMenu(m)
	fatal(app.Run())

	trayTemplate := []menu.Item{
		{
			Title: "Click on this here thing",
		},
		{
			Title: "Secret stuff",
			SubMenu: []menu.Item{
				{
					ID:    42,
					Title: "I'm nested!!",
				},
				{
					ID:       101,
					Title:    "Can't touch this",
					Disabled: true,
				},
			},
		},
		{
			Title:       "Quit App",
			Accelerator: "Command+T",
		},
	}

	iconPath := "assets/icon.png"
	if runtime.GOOS == "windows" {
		iconPath = "assets/icon.ico"
	}

	iconData, err := ioutil.ReadFile(iconPath)
	if err != nil {
		fmt.Println("Error reading icon file:", err)
	}

	platform.Dispatch(func() {
		app.NewIndicator(iconData, trayTemplate)
	})

	options := window.Options{
		Title: "Demo window",
		// NOTE(nick): resizing a transparent window on MacOS seems really slow?
		Transparent: false,
		Frameless:   false,
		Visible:     true,
		Resizable:   true,
		//Position: window.Position{X: 10, Y: 10},
		Size:   window.Size{Width: 360, Height: 240},
		Center: true,
		HTML: `
			<!doctype html>
			<html>
				<body style="font-family: -apple-system, BlinkMacSystemFont, avenir next, avenir, segoe ui, helvetica neue, helvetica, Ubuntu, roboto, noto, arial, sans-serif; background-color:rgba(87,87,87,0.8);"></body>
				<script>
					window.onload = function() {
						document.body.innerHTML = '<div style="padding: 30px">Transparency Test!<br><br>${navigator.userAgent}</div>';
					};
				</script>
			</html>
		`,
	}

	core.Dispatch(func() {
		w1, err := window.New(options)
		fatal(err)

		fmt.Println("[main] window", w1)

		w1.SetTitle("Hello, Sailor!")
		fmt.Println("[main] window position", w1.GetOuterPosition())
	})

	// shell.ShowNotification(shell.Notification{
	// 	Title:    "Title: Hello, world",
	// 	Subtitle: "Subtitle: MacOS only",
	// 	Body:     "Body: This is the body",
	// })

	// ok := shell.ShowMessage(shell.MessageDialog{
	// 	Title:   "Title: what do you think?",
	// 	Body:    "Body: about this description text",
	// 	Level:   "warning",
	// 	Buttons: "okcancel",
	// })

	// fmt.Println("ShowMessage ok", ok)

	// files := shell.ShowFilePicker(shell.FileDialog{
	// 	Title:   "Title: please pick a file...",
	// 	Mode:    "pickfiles",
	// 	Filters: []string{"txt,rs,cpp", "image:png,jpg,jpeg"},
	// })

	// fmt.Println("ShowFilePicker files", files, len(files))

	// success := shell.WriteClipboard("Hello from Go!")
	// fmt.Println("Wrote clipboard data:", success)

	// fmt.Println("Read clipboard data:", shell.ReadClipboard())

	core.Dispatch(func() {
		displays := system.Displays()
		fmt.Println("Displays:")

		for _, it := range displays {
			fmt.Println("", it.Name)
			fmt.Println("  Size:", it.Size)
			fmt.Println("  Position:", it.Position)
			fmt.Println("  ScaleFactor:", it.ScaleFactor)
		}
	})

	// didRegister1 := shell.RegisterShortcut("Control+Shift+R")
	// fmt.Println("didRegister", didRegister1)

	// didRegister2 := shell.RegisterShortcut("Control+Shift+T")
	// fmt.Println("didRegister", didRegister2)

	// didUnregister := shell.UnregisterShortcut("Control+Shift+T")
	// fmt.Println("didUnregister", didUnregister)

}