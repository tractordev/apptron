package main

import (
	"tractor.dev/hostbridge/app"
)

//go:embed *
var dir embed.FS

type rpc struct{}

func main() {
	app.Run(&rpc{}, dir)
}
