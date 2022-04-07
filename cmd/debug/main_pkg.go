//go:build pkg

package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/progrium/macdriver/core"
	"tractor.dev/hostbridge/bridge/api/app"
	"tractor.dev/hostbridge/bridge/api/menu"
	"tractor.dev/hostbridge/bridge/api/shell"
	"tractor.dev/hostbridge/bridge/api/system"
	"tractor.dev/hostbridge/bridge/api/window"
	"tractor.dev/hostbridge/bridge/event"
	"tractor.dev/hostbridge/bridge/misc"
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

	event.Listen(struct{}{}, func(e event.Event) error {
		log.Println(e)
		return nil
	})

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
	fatal(app.Run(app.Options{}))

	defer shell.UnregisterAllShortcuts()

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
			Title:       "Quit",
			Accelerator: "Command+T",
		},
	}

	iconData, err := misc.Assets.ReadFile("icon.png")
	fatal(err)

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

	shell.RegisterShortcut("CMD+SHIFT+S")

	core.Dispatch(func() {
		w1, err := window.New(options)
		fatal(err)

		fmt.Println("[main] window", w1)

		w1.SetTitle("Hello, Sailor!")
		fmt.Println("[main] window position", w1.GetOuterPosition())
	})

	core.Dispatch(func() {
		shell.ShowNotification(shell.Notification{
			Title:    "Title: Hello, world",
			Subtitle: "Subtitle: MacOS only",
			Body:     "Body: This is the body",
		})
	})

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

	core.Dispatch(func() {
		shell.WriteClipboard("Hello from Go!")
		fmt.Println("Read written clipboard data:", shell.ReadClipboard())
	})

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

	core.Dispatch(func() {
		if shell.ShowMessage(shell.MessageDialog{
			Title:   "TITLE",
			Level:   "error",
			Buttons: "okcancel",
			Body:    "BODY",
		}) {
			fmt.Println("YES")
		} else {
			fmt.Println("No")
		}
		shell.UnregisterShortcut("CMD+SHIFT+S")
	})

	// core.Dispatch(func() {
	// 	ret := shell.ShowFilePicker(shell.FileDialog{
	// 		Directory: "/Users/progrium/Source/github.com/tractordev/hostbridge",
	// 		Filters:   []string{"go,js"},
	// 		Title:     "TITLE",
	// 		Mode:      "pickfiles",
	// 	})
	// 	fmt.Println("file picker:", ret)
	// })

	// go func() {
	// 	<-time.After(3 * time.Second)
	// 	platform.Dispatch(func() {
	// 		mnu := menu.New([]menu.Item{
	// 			{
	// 				ID:    1,
	// 				Title: "Hello",
	// 			},
	// 			{
	// 				ID:    2,
	// 				Title: "One",
	// 			},
	// 			{
	// 				ID:    3,
	// 				Title: "Two",
	// 			},
	// 		})
	// 		mnu.Popup()
	// 	})
	// }()

	// didRegister1 := shell.RegisterShortcut("Control+Shift+R")
	// fmt.Println("didRegister", didRegister1)

	// didRegister2 := shell.RegisterShortcut("Control+Shift+T")
	// fmt.Println("didRegister", didRegister2)

	// didUnregister := shell.UnregisterShortcut("Control+Shift+T")
	// fmt.Println("didUnregister", didUnregister)

	select {}
}
