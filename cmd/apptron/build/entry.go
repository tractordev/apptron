//go:build exclude

package main

import (
	"embed"

	"tractor.dev/apptron/apputil"
	"tractor.dev/apptron/client"
)

//go:embed *
var dir embed.FS

type rpc struct {
	bridge *client.Client
}

func (r *rpc) SetClient(b *client.Client) {
	r.bridge = b
}

func main() {
	apputil.Run(dir, &rpc{})
}
