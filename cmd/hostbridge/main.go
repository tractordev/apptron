package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/progrium/qtalk-go/mux"
	"tractor.dev/hostbridge/bridge"
	"tractor.dev/hostbridge/bridge/platform"
	"tractor.dev/hostbridge/cmd/hostbridge/build"
)

const Version = "0.2.0"

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

	if *flagDebug {
		fmt.Fprintf(os.Stderr, "hostbridge %s\n", Version)
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
