package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/progrium/hostbridge/bridge"
	"github.com/progrium/qtalk-go/mux"
)

const Version = "0.1.0"

func main() {
	flagDebug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	if *flagDebug {
		fmt.Fprintf(os.Stderr, "hostbridge %s\n", Version)
	}

	sess, err := mux.DialIO(os.Stdout, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	srv := bridge.NewServer()
	srv.Respond(sess, context.Background())
}
