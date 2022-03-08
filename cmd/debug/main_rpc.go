//go:build rpc

package main

import (
	"context"
	"fmt"
	"net"
	"runtime"

	"github.com/progrium/hostbridge/bridge"
	"github.com/progrium/hostbridge/bridge/core"
	"github.com/progrium/hostbridge/client"
)

var quitId uint16 = 999
var quitAllId uint16 = 9999

func init() {
	runtime.LockOSThread()
}

func main() {
	go Run()
	core.Run(nil)
}

func Run() {
	// setup server
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	srv := bridge.NewServer()
	go srv.Serve(l)

	// setup client
	c, err := client.Dial(l.Addr().String())
	if err != nil {
		panic(err)
	}
	fmt.Println(c)
	//defer c.Close()

	ctx := context.Background()

	menuTemplate := []client.MenuItem{
		{
			// NOTE(nick): when setting the window menu with wry, the first item title will always be the name of the executable on MacOS
			// so, this property is ignored:
			// @Robustness: maybe we want to make that more visible to the user somehow?
			Title:   "this doesnt matter",
			Enabled: true,
			SubMenu: []client.MenuItem{
				{
					ID:          121,
					Title:       "About",
					Enabled:     true,
					Accelerator: "Control+I",
				},
				{
					ID:      122,
					Title:   "Disabled",
					Enabled: false,
				},
				{
					ID:          quitId,
					Title:       "Quit",
					Enabled:     true,
					Accelerator: "CommandOrControl+Q",
				},
			},
		},
		{
			ID:      23,
			Title:   "hello world",
			Enabled: true,
			SubMenu: []client.MenuItem{
				{
					ID:      777,
					Title:   "This is an amazing menu option",
					Enabled: true,
				},
			},
		},
	}

	m, err := c.Menu.New(ctx, menuTemplate)
	if err != nil {
		panic(err)
	}

	fmt.Println("Menu", m)
	if err = c.App.SetMenu(ctx, m); err != nil {
		panic(err)
	}

	/*
	trayTemplate := []client.MenuItem{
		{
			Title:   "Click on this here thing",
			Enabled: true,
		},
		{
			Title:   "Secret stuff",
			Enabled: true,
			SubMenu: []client.MenuItem{
				{
					ID:      42,
					Title:   "I'm nested!!",
					Enabled: true,
				},
				{
					ID:      101,
					Title:   "Can't touch this",
					Enabled: false,
				},
			},
		},
		{
			ID:          quitAllId,
			Title:       "Quit App",
			Enabled:     true,
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

	if err = c.App.NewIndicator(ctx, iconData, trayTemplate); err != nil {
		panic(err)
	}
	*/

	options := client.WindowOptions{
		Title: "Demo window",
		// NOTE(nick): resizing a transparent window on MacOS seems really slow?
		Transparent: true,
		Frameless:   false,
		Visible:     true,
		//Position: window.Position{X: 10, Y: 10},
		//Size: window.Size{ Width: 360, Height: 240 },
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
	if p, err := w1.GetOuterPosition(ctx); err != nil {
		panic(err)
	} else {
		fmt.Println("[main] window position", p)
	}

	err = c.Shell.ShowNotification(ctx, client.Notification{
		Title:    "Title: Hello, world",
		Subtitle: "Subtitle: MacOS only",
		Body:     "Body: This is the body",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("[main] Run done")
}
