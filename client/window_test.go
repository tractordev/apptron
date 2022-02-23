package client

import (
	"context"
	"fmt"
<<<<<<< HEAD
	"io/ioutil"
	"path"
	"runtime"
=======
>>>>>>> 6bff410 (refactor bridge packages to export to rpc easily, start client with tests)
	"testing"
)

func TestWindowModule(t *testing.T) {
	client, cleanup := setupBridgeClient(t)
	defer cleanup()

	ctx := context.Background()

	_, filename, _, _ := runtime.Caller(0)
	iconpath := path.Join(path.Dir(filename), "../assets/icon.png")
	icon, err := ioutil.ReadFile(iconpath)
	if err != nil {
		t.Fatal(err)
	}

	opts := WindowOptions{
		Visible: true,
		Icon:    icon,
		HTML: `
			<!doctype html>
			<html>
				<body style="font-family: -apple-system, BlinkMacSystemFont, avenir next, avenir, segoe ui, helvetica neue, helvetica, Ubuntu, roboto, noto, arial, sans-serif; background-color:rgba(87,87,87,0.8);"></body>
				<script>
					window.onload = function() {
						document.body.innerHTML = '<div style="padding: 30px">TEST</div>';
					};
				</script>
			</html>
		`,
	}
	w1, err := client.Window.New(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	w1.OnMoved = func(e Event) {
		fmt.Println("MOVED:", e.Position)
	}
	w1.OnResized = func(e Event) {
		fmt.Println("RESIZED:", e.Size)
	}
	w1.SetSize(ctx, Size{Width: 1240, Height: 480})
	w1.SetTitle(ctx, "WINDOW1")
	w1.SetPosition(ctx, Position{X: 200, Y: 50})

	w2, err := client.Window.New(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	w2.SetSize(ctx, Size{Width: 640, Height: 48})
	w2.SetTitle(ctx, "WINDOW2")
	w2.SetPosition(ctx, Position{X: 200, Y: 100})

	w1.Focus(ctx)

	fmt.Println(w1, w2)
	// time.Sleep(5 * time.Second)
<<<<<<< HEAD

=======
>>>>>>> 6bff410 (refactor bridge packages to export to rpc easily, start client with tests)
}
