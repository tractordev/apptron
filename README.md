# Apptron (early access)

Apptron gives you webview windows and common platform APIs for your simple scripts, homebrew utilities, or full applications. Building cross-platform (Win, Mac, Linux) programs that leverage native functionality like menus, dialogs, notifications, global shortcuts, and others, has never been more accessible.

The Apptron executable provides these cross-platform API modules:

- **window**: create and manage native windows with a native (non-Chromium) webview
- **menu**: create and manage native menus
- **app**: manage desktop application properties (lifecycle, default menus, startup mode, icons, etc)
- **shell**: native desktop shell dialogs, notifications, clipboard, global shortcuts
- **system**: system resource information (displays, cpu/memory, OS info)

These modules can be used a number of ways:

- **API**: run as a subprocess and use full API over STDIO from any language using a library
- **CLI**: shell script friendly command line versions of most of the API
- **build**: build a webview binary with `apptron build` and use the API from pages in JavaScript
- **import**: use module packages directly in Go (forces CGO, takes up main thread, advanced use only)

## Getting Apptron

TODO

## Using Apptron CLI

Much of the Apptron API can be used from the command-line, letting you use it in shell scripts.

```
$ apptron -h
```

## Using Apptron API

TODO

## Getting Help

There is a `#apptron` channel in the Progrium Discord. Feel free to ask for help or get involved in development there or file issues or PRs here.

## Contributing

TODO

## Developing

TODO

The dependencies for building are Go, and platform libraries (for example, XCode Developer Tools on Mac).

## License

MIT