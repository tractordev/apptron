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
	"text/tabwriter"

	"github.com/progrium/qtalk-go/mux"
	"tractor.dev/apptron/bridge"
	"tractor.dev/apptron/bridge/api/shell"
	"tractor.dev/apptron/bridge/api/system"
	"tractor.dev/apptron/bridge/platform"
	"tractor.dev/apptron/cmd/apptron/build"
	"tractor.dev/apptron/cmd/apptron/bundle"
)

const Version = "0.3.0"

func init() {
	runtime.LockOSThread()
}

func main() {
	flagDebug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	if flag.Arg(0) == "build" {
		build.Build()
		return
	}

	if flag.Arg(0) == "clean" {
		build.Clean()
		return
	}

	if flag.Arg(0) == "bundle" {
		bundle.Bundle()
		return
	}

	if flag.Arg(0) == "app" && flag.Arg(1) == "indicator" {
		appIndicator()
		return
	}

	if flag.Arg(0) == "menu" && flag.Arg(1) == "popup" {
		menuPopup()
		return
	}

	if flag.Arg(0) == "system" && flag.Arg(1) == "displays" {
		systemDisplays()
		return
	}

	if flag.Arg(0) == "shell" && flag.Arg(1) == "notification" {
		shellNotification()
		return
	}

	if flag.Arg(0) == "shell" && flag.Arg(1) == "message" {
		shellMessage()
		return

	}
	if flag.Arg(0) == "shell" && flag.Arg(1) == "file-picker" {
		shellFilePicker()
		return
	}

	if flag.Arg(0) == "shell" && flag.Arg(1) == "clipboard" {
		shellClipboard()
		return
	}

	if flag.Arg(0) == "shell" && flag.Arg(1) == "shortcuts" {
		shellShortcuts()
		return
	}

	if *flagDebug {
		fmt.Fprintf(os.Stderr, "apptron %s\n", Version)
	}

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
}

//apptron app indicator <menu-path> <icon-path>
func appIndicator() {

}

//apptron menu popup <menu-path>
func menuPopup() {

}

//apptron system displays
func systemDisplays() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, d := range system.Module.Displays() {
		fmt.Fprintf(w, "Name:\t%s\n", d.Name)
		fmt.Fprintf(w, "Size:\t%.fx%.f\n", d.Size.Width, d.Size.Height)
		fmt.Fprintf(w, "Scale:\t%.f\n", d.ScaleFactor)
		fmt.Fprintf(w, "Position:\t%.fx%.f\n", d.Position.X, d.Position.Y)
		fmt.Fprintln(w, "")
	}
	w.Flush()
}

//apptron shell notification <title> <body>
func shellNotification() {

}

//apptron shell message [--buttons=yesno,..] [--level=info,..] <title> <body>
func shellMessage() {

}

//apptron shell file-picker <mode> [<title>] [--directory=...] [--filepath=...] [--filters=text:dat,txt;source:go,html]
func shellFilePicker() {

}

//apptron shell clipboard [write-string]
func shellClipboard() {
	write := flag.Arg(2)
	if write != "" {
		if write == "-" {
			b, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
			write = string(b)
		}
		if !shell.Module.WriteClipboard(write) {
			os.Exit(1)
		}
		return
	}
	fmt.Println(shell.Module.ReadClipboard())
}

//apptron shell shortcuts <accelerator[:program]>...
func shellShortcuts() {

}
