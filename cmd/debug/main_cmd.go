//go:build cmd

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tractordev/hostbridge/client"
)

func main() {
	os.Setenv("BRIDGECMD", "./hostbridge")
	c, err := client.Spawn()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	ctx := context.Background()

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

	select {}
}
