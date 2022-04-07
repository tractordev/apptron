//go:build app

package main

import (
	"embed"
	"os"

	"tractor.dev/hostbridge/apputil"
	"tractor.dev/hostbridge/client"
)

//go:embed index.html icon.png
var dir embed.FS

type rpc struct {
	bridge *client.Client
}

func (r *rpc) SetClient(b *client.Client) {
	r.bridge = b
}

func main() {
	os.Setenv("BRIDGECMD", "./hostbridge")
	apputil.Run(dir, &rpc{})
}
