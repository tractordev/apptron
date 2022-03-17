package main

import (
	"os"

	"tractor.dev/hostbridge/app"
)

type rpc struct{}

func main() {
	app.Run(&rpc{})
}
