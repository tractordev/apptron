package main

import (
	"context"
	"log"
	"os"

	"tractor.dev/toolkit-go/engine/cli"
)

var Version = "dev"

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	root := &cli.Command{
		Version: Version,
		Usage:   "aptn",
		Long:    `aptn provides utilities for Apptron environments`,
	}

	root.AddCommand(execCmd())
	root.AddCommand(fuseCmd())
	root.AddCommand(portsCmd())
	root.AddCommand(shm9pCmd())
	root.AddCommand(shmtestCmd())

	if err := cli.Execute(context.Background(), root, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
