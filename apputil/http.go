package apputil

import (
	"net/http"

	_ "embed"

	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/mux"
	"github.com/progrium/qtalk-go/rpc"
	"golang.org/x/net/websocket"
	"tractor.dev/hostbridge/client"
	"tractor.dev/hostbridge/clientjs/dist"
)

//go:embed loader.js
var loader []byte

// BackendServer returns an http.Handler that responds to three endpoints,
// all hardcoded with /-/ path prefix:
//
// 	/-/ws: 						the WebSocket endpoint
// 	/-/client.js: 		the JavaScript client module
// 	/-/hostbridge.js:	the JavaScript loader
//
// The WebSocket endpoint establishes a qtalk session that will proxy to the provided
// hostbridge Client. The client module served is the embeded JS client
// from the clientjs/dist "package". The loader is what would be included by an HTML
// page to create a $host global that points to the hostbridge client instance
// created after connecting to the WebSocket endpoint. The $host global also has a
// ready property which is a Promise resolved when the client is connected and ready
// to be used.
//
// The optional muxExt callback allows you to handle custom RPC selectors that will
// be matched against before attempting to pass them to the hostbridge Client.
//
// If none of the endpoints are matched against, a NotFound error is returned. If
// serving other endpoints, put this behind a check for the /-/ path prefix.
func BackendServer(proxyTo *client.Client, muxExt func(mux *rpc.RespondMux)) http.Handler {
	return &backend{proxyTo, muxExt}
}

type backend struct {
	proxyTo *client.Client
	muxExt  func(mux *rpc.RespondMux)
}

func (b *backend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/-/client.js":
		w.Header().Add("Content-Type", "text/javascript")
		w.Write(dist.ClientJS)
	case "/-/hostbridge.js":
		w.Header().Add("Content-Type", "text/javascript")
		w.Write(loader)
	case "/-/ws":
		websocket.Handler(func(ws *websocket.Conn) {
			ws.PayloadType = websocket.BinaryFrame
			sess := mux.New(ws)

			// recreate a default handler on proxyTo
			// that proxies to this session (for callbacks)
			proxyToMux := b.proxyTo.Peer.RespondMux
			proxyToMux.Remove("")
			proxyToMux.Handle("", rpc.ProxyHandler(rpc.NewClient(sess, codec.JSONCodec{})))

			// create a server to proxy calls to proxyTo
			handler := rpc.NewRespondMux()
			handler.Handle("", rpc.ProxyHandler(b.proxyTo.Client))
			if b.muxExt != nil {
				b.muxExt(handler)
			}
			srv := &rpc.Server{
				Codec:   codec.JSONCodec{},
				Handler: handler,
			}
			srv.Respond(sess, nil)
		}).ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}
