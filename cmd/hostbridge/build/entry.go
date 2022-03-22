//go:build exclude

package main

import (
	"embed"

	"tractor.dev/hostbridge/apputil"
)

//go:embed *
var dir embed.FS

type rpc struct{}

func main() {
	apputil.Run(dir, &rpc{})
}
