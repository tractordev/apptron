// apptron is the main command of apptron
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/progrium/qtalk-go/mux"
	"golang.design/x/hotkey/mainthread"
	"tractor.dev/apptron/bridge"
	"tractor.dev/apptron/bridge/api/app"
	"tractor.dev/apptron/bridge/api/menu"
	"tractor.dev/apptron/bridge/api/shell"
	"tractor.dev/apptron/bridge/api/system"
	"tractor.dev/apptron/bridge/event"
	"tractor.dev/apptron/bridge/misc"
	"tractor.dev/apptron/bridge/platform"
	"tractor.dev/apptron/cmd/apptron/build"
	"tractor.dev/apptron/cmd/apptron/bundle"
	"tractor.dev/apptron/cmd/apptron/cli"
)

const Version = "0.3.0"

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
		Usage: "dev",
		Short: "launch a webview window from HTML",
		Run: func(ctx context.Context, args []string) {
			//build.Build()
		},
	})

	root.AddCommand(&cli.Command{
		Usage: "build",
		Short: "compile webview app from HTML",
		Run: func(ctx context.Context, args []string) {
			build.Build()
		},
	})

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

func windowLaunch() *cli.Command {
	return &cli.Command{
		Usage: "launch",
		Run: func(ctx context.Context, args []string) {
			// TODO
		},
	}
}

func appIndicator() *cli.Command {
	return &cli.Command{
		Usage: "indicator <menu-path> [icon-path]",
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
		Args:  cli.ExactArgs(2),
		Run: func(ctx context.Context, args []string) {
			runApp(func() {
				shell.Module.ShowNotification(shell.Notification{
					Title: flag.Arg(2),
					Body:  flag.Arg(3),
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
		Usage: "file-picker [flags...] <mode> [title]", //pickfile, pickfiles, pickfolder, savefile
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
			runApp(func() {
				paths := shell.Module.ShowFilePicker(shell.FileDialog{
					Mode:      args[0],
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
	cmd.Flags().StringVar(&directory, "directory", "", "directory to start in")
	cmd.Flags().StringVar(&filename, "filename", "", "default file selected")
	cmd.Flags().StringVar(&filters, "filters", "", "extension filters [ex: txt or text:dat,txt;source:go,html]")
	return cmd
}

func shellClipboard() *cli.Command {
	return &cli.Command{
		Usage: "clipboard [string]",
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
