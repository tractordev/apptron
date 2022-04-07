# apptron (alpha)

*This is currently sponsorware, so it is only available to sponsors or others with explicit access. You may run into toolchain issues working with a private repository, if so please refer to this [wiki page](https://github.com/tractordev/tractordev.github.io/wiki/Private-Repository).*

Apptron is an executable that exposes common APIs for building webview-based desktop applications and integrating with the GUI shell in a way that is platform (Win, Mac, Linux) and language agnostic. This makes it a simple, compact primitive for building your own native cross-platform Electron-style applications, alternative or specialized Electron-style frameworks, or even simple scripts for customizing your desktop experience. 

*NOTE: Windows and Linux support is only availble in the v0.1 tag while re-implementing the bridge backends*

It provides these high-level API modules:

- **window**: create and manage native windows with a native (non-Chromium) webview
- **menu**: create and manage native menus
- **app**: manage desktop application properties (lifecycle, default menus, startup mode, icons, etc)
- **shell**: native desktop shell dialogs, notifications, clipboard, global shortcuts
- **system**: system resource information (displays, cpu/memory, OS info)

Apptron runs as a subprocess helper with communication over STDIO using [qtalk](https://github.com/tractordev/qtalk). Running as a subprocess means regardless of language it is always distributed as the same single binary, frees your program from the restrictive and error-prone GUI app threading model, and avoids native extensions/FFI that can complicate the build and install for your program.

Currently there are two languages supported: Go and JavaScript. Adding language support is straightforward and can be done using [this guide]().

## Getting and distributing

A framework, application, or tool built with apptron may choose to embed the apptron binary within itself, but soon it will be available to download as a standalone binary for use as a package dependency or to write scripts against directly. For now, you can clone this repo and build the binary with:

```bash
make apptron
```

The dependencies for building are [Go](https://go.dev/), [Rust](https://www.rust-lang.org/tools/install) (for the moment), and platform libraries (for example, XCode Developer Tools on Mac).

## Builtin Electron alternative

Since apptron is just an API for native app functionality, you still need to write and bundle your application. This is why apptron can be seen as a runtime primitive for building your own Electron framework. To get started quickly, there is a build subcommand that can take some web assets and produce a single distributable binary. This build command is built-in to the apptron binary, but also requires Go installed to run. Perhaps eventually Go can be bundled with it as well.

This framework is very simple, much simpler than Electron. The minimal workflow is:

 - build or install `apptron`
 - make an `index.html` file 
 - run `apptron build` in that directory
 - get an executable that contains and runs the page in a native window
 
Any other assets in the directory are bundled and served as well. Any page can include `/-/apptron.js` and then use the `$host` global to access the apptron API to create other windows, etc letting you build an entire native application with HTML and JavaScript, without a drop of NodeJS or NPM modules.

From the `demo` directory, which includes an example HTML file for this, assuming `apptron` was built in the project root, you can run:

```bash
../apptron build
```

There should now be a `demo` executable (named after the directory) that starts a native application based on your web assets. The main window can be customized from the HTML by including a meta tag in the `index.html` like `<meta name="window" content="center=true,width=800,height=600" />`.

You can even extend the capabilities accessible from `$host` by including a Go source file in the directory that add methods to an `rpc` struct type of the `main` package. Any method on this struct will be callable from JavaScript.

Since this framework and builder is meant as an *example* and a quick development tool, it doesn't, for example, create a signed app bundle with manifest, etc. Only an executable, which you can then either manually or using another tool make into an app bundle fit for an App Store. 

## Building your own application

As long as you have a apptron client in your language, you can build a program against the apptron API giving you much more control over how you build your application, or whether it's even a full blown application. It could be a Python script. 

The above builder is basically just a wrapper around a pre-made Go program, which conveniently has support to bundle static files into the binary it produces. Rust has this capability built-in as well, but many languages have a tool or mechanism to do something like this.

To see how you can work against the apptron API directly, lets recreate a program similar to what the above builder wrote for us. This could be any language, but for convenience, we'll be using Go.

```golang
package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"strings"

	"tractor.dev/apptron/client"
	"tractor.dev/apptron/apputil"
)

//go:embed *.html
var assets embed.FS

func main() {
	// this starts the apptron subprocess
	// and returns a client for its API
	bridge, err := client.Spawn()
	if err != nil {
		log.Fatal(err)
	}

	bridge.OnEvent = func(event client.Event) {
		// just so we can see what kind of events
		// we get. most are from window interactions
		log.Println(event)
	}

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// handle requests for the /-/apptron.js file and
		// the WebSocket endpoint it connects to, which then
		// proxies to our apptron client
		if strings.HasPrefix(r.URL.Path, "/-/") {
			apputil.BackendServer(bridge, nil).ServeHTTP(w, r)
			return
		}
		// otherwise serve the embedded assets picked up at build time
		http.FileServer(http.FS(assets)).ServeHTTP(w, r)
	}))
	go http.ListenAndServe(":8888", nil)

	ctx := context.Background()

	// this starts the application event loop
	if err := bridge.App.Run(ctx, client.AppOptions{}); err != nil {
		log.Fatal(err)
	}

	// set up some default window options that can be overwritten by a
	// meta tag with name="window" in the index.html of our embedded assets
	opts := apputil.OptionsFromHTML(assets, "index.html", "window", client.WindowOptions{
		Size: client.Size{
			Width:  640,
			Height: 480,
		},
		Visible: true,
		URL:     "http://localhost:8888/",
	})

	// we create a window with our options
	if _, err = bridge.Window.New(ctx, opts); err != nil {
		log.Fatal(err)
	}

	// meanwhile wait until apptron exits,
	// most likely from our main window being closed
	bridge.Wait()
}

```

Now we can put this file in a directory with an `index.html` (maybe from the `demo` dir again) and simply run `go build` to get a binary.

Since we're doing similar things as our builder does, we can reuse some of those utility functions. However, you could also make a program without them, and instead doing all apptron calls from Go, perhaps exposing them to JavaScript some other way or not at all.

For example, here is an even simpler Go program that doesn't even make a window (so it doesn't need to serve any files). It simply lets you pick some files using the native file picker dialog and then prints them out. 

```golang
package main

import (
	"context"
	"fmt"
	"log"

	"tractor.dev/apptron/client"
)

func main() {
	bridge, err := client.Spawn()
	if err != nil {
		log.Fatal(err)
	}
	defer bridge.Close()

	ctx := context.Background()
	files, err := bridge.Shell.ShowFilePicker(ctx, client.FileDialog{
		Mode: "pickfiles",
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file)
	}

}
```

For something as simple as using this one native dialog, we don't have to invoke CGO and compile against tons of system headers. apptron makes these APIs much more accessible for simple programs and scripts. 

## Getting help

There is a `#apptron` channel in the Progrium Discord. Feel free to ask for help or get involved in development there or file issues or PRs here.

## License

MIT