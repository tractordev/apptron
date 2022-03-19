//go:build exclude

package main

import (
	"embed"

	"tractor.dev/hostbridge/cmd/hostbridge/demo"
)

//go:embed *
var dir embed.FS

type rpc struct{}

func main() {
	demo.Run(&rpc{}, dir)
}
