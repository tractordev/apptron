package apputil

import (
	"net/http"
	"strings"

	_ "embed"

	"github.com/progrium/qtalk-go/mux"
	"github.com/progrium/qtalk-go/rpc"
	"github.com/progrium/qtalk-go/x/cbor/codec"
	"golang.org/x/net/websocket"
	"tractor.dev/apptron/chrome"
	"tractor.dev/apptron/client"
	"tractor.dev/apptron/clientjs/dist"
)

//go:embed loader.js
var loader []byte

// BackendServer returns an http.Handler that responds to builtin endpoints,
// all hardcoded with /-/ path prefix:
//
//	/-/ws: 						the WebSocket endpoint
//	/-/client.js: 		the JavaScript client module
//	/-/apptron.js:		the JavaScript loader
//	/-/chrome:				the builtin chrome pages dir
//
// The WebSocket endpoint establishes a qtalk session that will proxy to the provided
// apptron Client. The client module served is the embeded JS client
// from the clientjs/dist "package". The loader is what would be included by an HTML
// page to create a $host global that points to the apptron client instance
// created after connecting to the WebSocket endpoint. The $host global also has a
// ready property which is a Promise resolved when the client is connected and ready
// to be used.
//
// The optional muxExt callback allows you to handle custom RPC selectors that will
// be matched against before attempting to pass them to the apptron Client.
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
	case "/-/apptron.js":
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
			proxyToMux.Handle("", rpc.ProxyHandler(rpc.NewClient(sess, codec.CBORCodec{})))

			// create a server to proxy calls to proxyTo
			handler := rpc.NewRespondMux()
			handler.Handle("", rpc.ProxyHandler(b.proxyTo.Client))
			if b.muxExt != nil {
				b.muxExt(handler)
			}
			srv := &rpc.Server{
				Codec:   codec.CBORCodec{},
				Handler: handler,
			}
			srv.Respond(sess, nil)
		}).ServeHTTP(w, r)
	default:
		if strings.HasPrefix(r.URL.Path, "/-/chrome") {
			http.StripPrefix("/-/chrome", http.FileServer(http.FS(chrome.Dir))).ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	}
}
