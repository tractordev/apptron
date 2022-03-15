package main

import (
	"os"

	"github.com/tractordev/hostbridge/client"
)

func main() {
	// run hostbridge
	os.Setenv("BRIDGECMD", "./hostbridge")
	c, err := client.Spawn()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// serve local dir over http

	// accept qtalk ws connections, pass channels to hostbridge
	// read window options from index.html, show window
}
