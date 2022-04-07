//go:build app

package main

import (
	"embed"
	"os"

	"tractor.dev/apptron/apputil"
	"tractor.dev/apptron/client"
)

//go:embed index.html
var dir embed.FS

type rpc struct {
	bridge *client.Client
}

func (r *rpc) SetClient(b *client.Client) {
	r.bridge = b
}

func main() {
	os.Setenv("BRIDGECMD", "./apptron")
	apputil.Run(dir, &rpc{})
}
