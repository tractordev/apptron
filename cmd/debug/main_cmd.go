//go:build cmd

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"tractor.dev/hostbridge/client"
)

func main() {
	os.Setenv("BRIDGECMD", "./hostbridge")
	c, err := client.Spawn()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c.OnEvent = func(e client.Event) {
		log.Println(e)
	}

	ctx := context.Background()

	if err := c.App.Run(ctx); err != nil {
		panic(err)
	}

	options := client.WindowOptions{
		Title: "Demo window",
		// NOTE(nick): resizing a transparent window on MacOS seems really slow?
		Transparent: true,
		Frameless:   false,
		Resizable:   true,
		Visible:     true,
		//Position: window.Position{X: 10, Y: 10},
		Size:   client.Size{Width: 360, Height: 240},
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

	w1, err := c.Window.New(ctx, options)
	if err != nil {
		panic(err)
	}

	fmt.Println("[main] window", w1)

	if w1 == nil {
		return
	}

	if err := w1.SetTitle(ctx, "Hello, Sailor!"); err != nil {
		panic(err)
	}

	trayTemplate := []client.MenuItem{
		{
			Title: "Click on this here thing",
		},
		{
			Title: "Secret stuff",
			SubMenu: []client.MenuItem{
				{
					ID:    1,
					Title: "I'm nested!!",
				},
				{
					ID:       101,
					Disabled: true,
					Title:    "Can't touch this",
				},
			},
		},
		{
			Title:       "Quit",
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

	if err := c.App.NewIndicator(ctx, iconData, trayTemplate); err != nil {
		log.Fatal(err)
	}

	// m, err := c.Menu.New(ctx, []client.MenuItem{
	// 	{
	// 		ID:      10,
	// 		Title:   "AAA",
	// 		Enabled: true,
	// 	},
	// 	{
	// 		ID:      20,
	// 		Title:   "BBB",
	// 		Enabled: true,
	// 	},
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err := c.Menu.Popup(ctx, m); err != nil {
	// 	log.Fatal(err)
	// }

	select {}
}
