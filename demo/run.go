//go:build !entrypoint

package demo

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/qtalk-go/mux"
	"github.com/progrium/qtalk-go/rpc"
	"golang.org/x/net/html"
	"golang.org/x/net/websocket"
	"tractor.dev/hostbridge/client"
	"tractor.dev/hostbridge/clientjs/dist"

	_ "embed"
)

func Run(delegate interface{}, fsys fs.FS) {
	_, err := exec.LookPath("hostbridge")
	if err != nil {
		f, _ := fsys.Open("hostbridge")
		d, _ := ioutil.ReadAll(f)
		f.Close()
		dir, err := ioutil.TempDir("", "hostbridge-*")
		if err != nil {
			log.Fatal(err)
		}
		path := filepath.Join(dir, "hostbridge")
		if err := ioutil.WriteFile(path, d, 0755); err != nil {
			log.Fatal(err)
		}
		os.Setenv("BRIDGECMD", path)
		log.Println(path)
	}

	c, err := client.Spawn()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	c.OnEvent = func(event client.Event) {
		if event.Type == client.EventClose && event.WindowID == 1 {
			c.Close()
		}
	}

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") != "websocket" {
			if r.URL.Path == "/-/client.js" {
				w.Header().Add("Content-Type", "text/javascript")
				w.Write(dist.ClientJS)
				return
			}
			http.FileServer(http.FS(fsys)).ServeHTTP(w, r)
			return
		}
		websocket.Handler(func(ws *websocket.Conn) {
			ws.PayloadType = websocket.BinaryFrame
			handler := rpc.NewRespondMux()
			handler.Handle("user", fn.HandlerFrom(delegate))
			handler.Handle("", rpc.ProxyHandler(c.Client))
			srv := &rpc.Server{
				Codec:   codec.JSONCodec{},
				Handler: handler,
			}
			srv.Respond(mux.New(ws), nil)
		}).ServeHTTP(w, r)
	}))
	go http.ListenAndServe(":7778", nil)

	params, err := extractWindowParams(fsys)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	ctx := context.Background()
	_, err = c.Window.New(ctx, client.WindowOptions{
		Title:       params["title"],
		AlwaysOnTop: parseBool(params["alwaysontop"], false),
		Fullscreen:  parseBool(params["fullscreen"], false),
		Maximized:   parseBool(params["maximized"], false),
		Resizable:   parseBool(params["resizable"], true),
		Transparent: parseBool(params["transparent"], false),
		Frameless:   parseBool(params["frameless"], false),
		Visible:     parseBool(params["visible"], true),
		Center:      parseBool(params["center"], false),
		Position: client.Position{
			X: parseFloat(params["x"], 0),
			Y: parseFloat(params["y"], 0),
		},
		Size: client.Size{
			Width:  parseFloat(params["width"], 640),
			Height: parseFloat(params["height"], 480),
		},
		MinSize: client.Size{
			Width:  parseFloat(params["min-width"], 0),
			Height: parseFloat(params["min-height"], 0),
		},
		MaxSize: client.Size{
			Width:  parseFloat(params["max-width"], 0),
			Height: parseFloat(params["max-height"], 0),
		},
		URL: "http://localhost:7778/",
	})
	if err != nil {
		log.Fatal(err)
	}

	c.Wait()
}

func parseBool(v string, fallback bool) bool {
	if v == "true" {
		return true
	}
	if v == "false" {
		return false
	}
	return fallback
}

func parseFloat(v string, fallback float64) float64 {
	f, err := strconv.ParseFloat(v, 64)
	if err == nil {
		return f
	}
	return fallback
}

func extractWindowParams(fsys fs.FS) (map[string]string, error) {
	f, err := fsys.Open("index.html")
	if err != nil {
		panic(err)
	}
	z := html.NewTokenizer(f)
	inTitle := false
	params := make(map[string]string)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return params, z.Err()
		case html.StartTagToken:
			tn, _ := z.TagName()
			if bytes.Equal(tn, []byte("body")) {
				return params, io.EOF
			}
			if bytes.Equal(tn, []byte("title")) {
				inTitle = true
			}
		case html.TextToken:
			if inTitle {
				params["title"] = string(z.Text())
				inTitle = false
			}
		case html.SelfClosingTagToken:
			tn, attr := z.TagName()
			if bytes.Equal(tn, []byte("meta")) && attr {
				var name, content string
				for attr {
					var k, v []byte
					k, v, attr = z.TagAttr()
					if bytes.Equal(k, []byte("name")) {
						name = string(v)
					}
					if bytes.Equal(k, []byte("content")) {
						content = string(v)
					}
				}
				if name == "main-window" {
					for _, part := range strings.Split(content, ",") {
						kv := strings.SplitN(part, "=", 2)
						params[kv[0]] = kv[1]
					}
				}
			}
		}
	}
}
