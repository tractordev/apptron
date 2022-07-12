// apptron is the main command of apptron
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/progrium/qtalk-go/mux"
	"github.com/progrium/qtalk-go/rpc"
	"golang.design/x/hotkey/mainthread"
	"tractor.dev/apptron"
	"tractor.dev/apptron/apputil"
	"tractor.dev/apptron/bridge"
	"tractor.dev/apptron/bridge/api/app"
	"tractor.dev/apptron/bridge/api/menu"
	"tractor.dev/apptron/bridge/api/shell"
	"tractor.dev/apptron/bridge/api/system"
	"tractor.dev/apptron/bridge/event"
	"tractor.dev/apptron/bridge/misc"
	"tractor.dev/apptron/bridge/platform"
	"tractor.dev/apptron/client"
	"tractor.dev/apptron/cmd/apptron/build"
	"tractor.dev/apptron/cmd/apptron/bundle"
	"tractor.dev/apptron/cmd/apptron/cli"
)

var Version = "dev"

func init() {
	runtime.LockOSThread()
}

func main() {
	var flagDebug bool
	root := &cli.Command{
		Version: Version,
		Usage:   "apptron",
		Long: `Apptron is a tool for scriptable native app functionality and webview
windows. Running without a subcommand starts the API service over STDIO.`,
		Run: func(ctx context.Context, args []string) {
			sess, err := mux.DialIO(os.Stdout, os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
			srv := bridge.NewServer()
			go srv.Respond(sess, context.Background())
			go func() {
				sess.Wait()
				platform.Terminate()
			}()
			platform.Main()
		},
	}
	root.Flags().BoolVar(&flagDebug, "debug", false, "debug mode")

	root.AddCommand(&cli.Command{
		Usage: "run",
		Short: "launch a webview window from HTML",
		Run: func(ctx context.Context, args []string) {
			url := ""
			if len(args) == 1 {
				url = args[0]
			}
			launchWindow(url)
		},
	})

	var flagSetup bool
	buildcmd := &cli.Command{
		Usage: "build",
		Short: "compile webview app from HTML",
		Run: func(ctx context.Context, args []string) {
			_, err := exec.LookPath("go")
			if err != nil {
				fmt.Println("Unable to find Go in your PATH.\n")
				fmt.Println("Use this URL to download Golang for your platform:")
				fmt.Println("  https://go.dev/doc/install\n")
				fmt.Println("EARLY ACCESS NOTE")
				fmt.Println("Make sure Git can access Apptron repo over SSH:")
				fmt.Println("  https://github.com/tractordev/tractordev.github.io/wiki/Private-Repository\n")
				os.Exit(1)
			}
			if flagSetup {
				// switch runtime.GOOS {
				// case "windows":
				// 	cmd := exec.Command(filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe"), "url.dll,FileProtocolHandler", url)
				// 	cmd.Run()
				// case "darwin":
				// 	cmd := exec.Command("open", url)
				// 	cmd.Run()
				// case "linux":
				// 	// TODO
				// default:
				// }
				return
			}
			build.Build(flagDebug)
		},
	}
	buildcmd.Flags().BoolVar(&flagDebug, "debug", false, "debug mode")
	buildcmd.Flags().BoolVar(&flagSetup, "setup", false, "setup compiler")
	root.AddCommand(buildcmd)

	root.AddCommand(&cli.Command{
		Usage: "clean",
		Short: "remove webview app build cache files",
		Run: func(ctx context.Context, args []string) {
			build.Clean()
		},
	})

	root.AddCommand(&cli.Command{
		Usage: "bundle",
		Short: "build platform application bundle",
		Run: func(ctx context.Context, args []string) {
			bundle.Bundle()
		},
	})

	app := &cli.Command{
		Usage: "app",
		Short: "app related API commands",
	}
	app.AddCommand(appIndicator())
	root.AddCommand(app)

	win := &cli.Command{
		Usage: "window",
		Short: "window related API commands",
	}
	win.AddCommand(windowLaunch())
	root.AddCommand(win)

	menu := &cli.Command{
		Usage: "menu",
		Short: "menu related API commands",
	}
	menu.AddCommand(menuPopup())
	root.AddCommand(menu)

	system := &cli.Command{
		Usage: "system",
		Short: "system related API commands",
	}
	system.AddCommand(systemDisplays())
	root.AddCommand(system)

	shell := &cli.Command{
		Usage: "shell",
		Short: "shell related API commands",
	}
	shell.AddCommand(shellNotification())
	shell.AddCommand(shellMessage())
	shell.AddCommand(shellFilePicker())
	shell.AddCommand(shellClipboard())
	shell.AddCommand(shellShortcuts())
	root.AddCommand(shell)

	if err := cli.Execute(context.Background(), root, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func runApp(fn func(), term bool) {
	go func() {
		if err := app.Module.Run(app.Options{
			Accessory: true,
			Agent:     true,
		}); err != nil {
			log.Fatal(err)
		}
		platform.Dispatch(fn)
		if term {
			platform.Terminate()
		} else {
			select {}
		}
	}()
	platform.Main()
}

// Item1
// Item2
// 	Subitem1
// 	  Subsubitem
// 	Subitem2
// ---
// Item3
func parseMenuFile(path string) (items []menu.Item, table map[int]string, err error) {
	var b []byte
	if path == "-" {
		b, err = ioutil.ReadAll(os.Stdin)
	} else {
		b, err = ioutil.ReadFile(path)
	}
	if err != nil {
		return nil, nil, err
	}
	table = make(map[int]string)
	lines := strings.Split(string(b), "\n")
	var itemStack [][]menu.Item
	depth := -1
	id := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := strings.Count(line, "\t")
		if indent < depth {
			for indent < depth {
				itemStack[depth-1][len(itemStack[depth-1])-1].SubMenu = itemStack[depth]
				itemStack = itemStack[:len(itemStack)-1]
				depth--
			}
		}
		if len(itemStack) == indent {
			itemStack = append(itemStack, []menu.Item{})
			depth++
		}

		id++
		item := menu.Item{ID: id, Title: strings.TrimSpace(line)}
		if strings.HasPrefix(item.Title, "-") {
			item.Title = ""
			item.Separator = true
		}
		table[id] = item.Title
		itemStack[indent] = append(itemStack[indent], item)

	}
	items = itemStack[0]
	return
}

func appIndicator() *cli.Command {
	return &cli.Command{
		Usage: "indicator <menu-path> [icon-path]",
		Short: "show app indicator with menu",
		Args:  cli.MinArgs(1),
		Run: func(ctx context.Context, args []string) {
			items, table, err := parseMenuFile(args[0])
			if err != nil {
				log.Fatal(err)
			}
			icon, err := misc.Assets.ReadFile("icon.png")
			if err != nil {
				log.Fatal(err)
			}
			if len(args) > 1 {
				icon, err = ioutil.ReadFile(args[1])
				if err != nil {
					log.Fatal(err)
				}
			}
			event.Listen("", func(e event.Event) error {
				if e.Type == event.MenuItem {
					fmt.Println(table[e.MenuItem])
				}
				return nil
			})
			runApp(func() {
				app.NewIndicator(icon, items)
			}, false)
		},
	}
}

func menuPopup() *cli.Command {
	return &cli.Command{
		Usage: "popup <menu-path>",
		Short: "popup a context menu",
		Args:  cli.ExactArgs(1),
		Run: func(ctx context.Context, args []string) {
			items, table, err := parseMenuFile(args[0])
			if err != nil {
				log.Fatal(err)
			}
			runApp(func() {
				id := menu.Module.Popup(items)
				fmt.Println(table[id])
			}, true)
		},
	}
}

func systemDisplays() *cli.Command {
	return &cli.Command{
		Usage: "displays",
		Short: "show display information",
		Run: func(ctx context.Context, args []string) {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			for _, d := range system.Module.Displays() {
				fmt.Fprintf(w, "Name:\t%s\n", d.Name)
				fmt.Fprintf(w, "Size:\t%.fx%.f\n", d.Size.Width, d.Size.Height)
				fmt.Fprintf(w, "Scale:\t%.f\n", d.ScaleFactor)
				fmt.Fprintf(w, "Position:\t%.fx%.f\n", d.Position.X, d.Position.Y)
				fmt.Fprintln(w, "")
			}
			w.Flush()
		},
	}
}

func shellNotification() *cli.Command {
	return &cli.Command{
		Usage: "notification <title> <body>",
		Short: "show desktop notification",
		Args:  cli.ExactArgs(2),
		Run: func(ctx context.Context, args []string) {
			runApp(func() {
				shell.Module.ShowNotification(shell.Notification{
					Title: args[0],
					Body:  args[1],
				})
			}, true)
		},
	}
}

func shellMessage() *cli.Command {
	var (
		buttons string
		level   string
	)
	cmd := &cli.Command{
		Usage: "message [flags...] <title> <body>",
		Short: "show message dialog",
		Args:  cli.ExactArgs(2),
		Run: func(ctx context.Context, args []string) {
			runApp(func() {
				if !shell.Module.ShowMessage(shell.MessageDialog{
					Title:   args[0],
					Body:    args[1],
					Buttons: buttons,
					Level:   level,
				}) {
					os.Exit(1)
				}
			}, true)
		},
	}
	cmd.Flags().StringVar(&buttons, "buttons", "ok", "buttons to show")
	cmd.Flags().StringVar(&level, "level", "info", "dialog level")
	return cmd
}

func shellFilePicker() *cli.Command {
	var (
		directory string
		filename  string
		filters   string
	)
	cmd := &cli.Command{
		Usage: "file-picker [flags...] <mode> [title]", //file, files, folder, savefile
		Short: "show file picker dialog",
		Args:  cli.MinArgs(1),
		Run: func(ctx context.Context, args []string) {
			title := ""
			if len(args) > 1 {
				title = args[1]
			}
			var fltrs []string
			if filters != "" {
				fltrs = strings.Split(filters, ";")
			}
			mode := args[0]
			if mode != "savefile" {
				// TODO: maybe push up to API?
				mode = fmt.Sprintf("pick%s", mode)
			}
			runApp(func() {
				paths := shell.Module.ShowFilePicker(shell.FileDialog{
					Mode:      mode,
					Title:     title,
					Directory: directory,
					Filename:  filename,
					Filters:   fltrs,
				})
				for _, p := range paths {
					fmt.Println(p)
				}
			}, true)
		},
	}
	cmd.Flags().StringVar(&directory, "dir", "", "directory to start in")
	cmd.Flags().StringVar(&filename, "default", "", "default file selected")
	cmd.Flags().StringVar(&filters, "filters", "", "extension filters [ex: txt or text:dat,txt;source:go,html]")
	return cmd
}

func shellClipboard() *cli.Command {
	return &cli.Command{
		Usage: "clipboard [string]",
		Short: "show or set the clipboard text value",
		Args:  cli.MaxArgs(1),
		Run: func(ctx context.Context, args []string) {
			if len(args) == 1 {
				if args[0] == "-" {
					b, err := ioutil.ReadAll(os.Stdin)
					if err != nil {
						log.Fatal(err)
					}
					args[0] = string(b)
				}
				if !shell.Module.WriteClipboard(args[0]) {
					os.Exit(1)
				}
				return
			}
			fmt.Println(shell.Module.ReadClipboard())
		},
	}
}

func shellShortcuts() *cli.Command {
	return &cli.Command{
		Usage: "shortcuts <accelerator>...",
		Short: "register global shortcuts",
		Run: func(ctx context.Context, args []string) {
			mainthread.Init(func() {
				for _, arg := range args {
					shell.RegisterShortcut(strings.ToLower(arg))
				}
				defer shell.UnregisterAllShortcuts()
				event.Listen("", func(e event.Event) error {
					if e.Type == event.Shortcut {
						fmt.Println(e.Shortcut)
					}
					return nil
				})
				select {}
			})
		},
	}
}

func windowLaunch() *cli.Command {
	return &cli.Command{
		Usage: "launch [URL]",
		Short: "launch a webview window from HTML",
		Run: func(ctx context.Context, args []string) {
			url := ""
			if len(args) == 1 {
				url = args[0]
			}
			launchWindow(url)
		},
	}
}

func launchWindow(url string) {
	// TODO: use url, live reloads, env var config

	ctx := context.Background()
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fsys := os.DirFS(dir)

	appOpts := apputil.AppOptionsFromHTML(fsys, "index.html", "application", client.AppOptions{})
	native, err := apptron.Run(ctx, appOpts)
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/-/", apputil.BackendServer(native, func(mux *rpc.RespondMux) {
		//mux.Handle("user", fn.HandlerFrom(&extensions{}))
	}))
	http.Handle("/", http.FileServer(http.FS(fsys)))
	srv := &http.Server{Handler: http.DefaultServeMux}
	go srv.Serve(l)

	_, err = native.Window.New(ctx, apputil.WindowOptionsFromHTML(fsys, "index.html", "window", client.WindowOptions{
		Size: client.Size{
			Width:  800,
			Height: 600,
		},
		Center:  true,
		Visible: true,
		URL:     fmt.Sprintf("http://%s/", l.Addr().String()),
	}))
	if err != nil {
		log.Fatal(err)
	}

	native.Wait()

}
