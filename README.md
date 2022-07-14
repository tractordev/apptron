# Apptron (early access)

Apptron gives you webview windows and common platform APIs for your simple scripts, homebrew utilities, or full applications. Building cross-platform (Win, Mac, Linux) programs that leverage native functionality (menus, dialogs, notifications, global shortcuts, etc) has never been more accessible.

The Apptron executable provides these cross-platform API modules:

- **window**: create and manage native windows with a native (non-Chromium) webview
- **menu**: create and manage native menus
- **app**: manage desktop application properties (lifecycle, default menus, startup mode, icons, etc)
- **shell**: native desktop shell dialogs, notifications, clipboard, global shortcuts
- **system**: system resource information (displays, cpu/memory, OS info)

*Early Access Note: Smaller APIs might not be implemented in Windows/Linux yet, please [report](https://github.com/tractordev/apptron/issues) anything not working on your platform.*

These modules can be used a number of ways:

- **API**: run as a subprocess and use full API over STDIO from any language using a library
- **CLI**: shell script friendly command line versions of most of the API
- **build**: build a webview binary with `apptron build` and use the API from pages in JavaScript
- **import**: use module packages directly in Go (forces CGO, takes up main thread, advanced use only)

## Getting Apptron

### Download Release

Still getting [release automation](https://github.com/tractordev/apptron/issues/44) up to speed, but there are [builds you can download](https://github.com/tractordev/apptron/releases). Let me know if you have problems with them or what kinds of warnings they present on your platform. The macOS build is signed and notarized, but the others aren't yet. Homebrew Tap coming soon.

### Build Source

I'd love you to build from source so we can iterate quickly. It's pretty painless. This means [installing Go](https://go.dev/dl/). Mac will need XCode Developer Tools installed. Linux has [more pre-requisites](https://github.com/tractordev/apptron/tree/main/bridge/platform/linux).

```bash
make apptron
```

This will produce an executable at `./dist/apptron` that you can put in your PATH.

## Using Apptron CLI

Besides `build` and a some other top level commands, much of the Apptron API can be used from the command-line, letting you use it in shell scripts. Commands should provide proper usage information via help. Please report anything missing.

```
$ apptron -h
Usage:
apptron [command]

Apptron is a tool for scriptable native app functionality and webview
windows. Running without a subcommand starts the API service over STDIO.

Available Commands:
  run              launch a webview window from HTML
  build            compile webview app from HTML
  clean            remove webview app build cache files
  bundle           build platform application bundle
  app              app related API commands
  window           window related API commands
  menu             menu related API commands
  system           system related API commands
  shell            shell related API commands

Flags:
  -debug
        debug mode
  -v
        show version

Use "apptron [command] -help" for more information about a command.
```

The `build` subcommand depends on Go, which shouldn't be an issue since we build from source right now. However, it pulls from this repository so you may run into issues related to working with private repositories. Refer to this [wiki page](https://github.com/tractordev/tractordev.github.io/wiki/Private-Repository) for resolving this.

## Using Apptron API

Apptron exposes a [qtalk](https://github.com/tractordev/qtalk) API over STDIO. We currently have a Go client that makes using this API easy. There is a JavaScript wrapper for once you have a channel established (over STDIO or sometimes WebSocket). Any other language will need at the very least a qtalk/qmux implementation. Submit or upvote an issue to prioritize support for your language if you can't contribute it yourself.

In the meantime, you can shell out and use the CLI commands, or you can use the API from HTML/JavaScript using `apptron build` or `apptron run`, or you can use Go.

### Using JavaScript+run

The `apptron run` command is similar to running `apptron build` and running the binary, but does not produce a binary. Try it in one of the `demos` directories and see how the API is used there.

### Using Go

Assuming you have an HTML file next to this, here is how you can build a webview app from Go. Notice you have to serve the contents of the webview yourself, but is trivial in Go:

```golang
package main

import (
	"context"
	"embed"
	"log"
	"net/http"

	"tractor.dev/apptron"
)

//go:embed *.html
var assets embed.FS

func main() {
	ctx := context.Background()

	native, err := apptron.Run(ctx, apptron.AppOptions{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = native.Window.New(ctx, apptron.WindowOptions{
		Visible: true,
		Size:    apptron.Size{Width: 800, Height: 600},
		Center:  true,
		Title:   "Demo Title",
		URL:     "http://localhost:9090",
	})
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(assets)))
	go http.ListenAndServe(":9090", nil)

	native.Wait()
}

```

## Getting Help

There is a `#apptron` channel in the [Progrium Discord](https://discord.gg/4zp9WMUtTD). Feel free to ask for help.

## Contributing

PLEASE help make this project ready to release. Contribute by submitting issues or especially PRs. Ask for help in Discord.

Since this repository is private, GitHub won't allow you to fork. With early access you do have write permission, so you can push to a branch with your username and submit a PR from your branch into main. 

## License

MIT
