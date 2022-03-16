package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/progrium/qtalk-go/mux"
	"github.com/tractordev/hostbridge/client"
	"golang.org/x/net/websocket"
)

func main() {
	os.Setenv("BRIDGECMD", "./hostbridge")
	c, err := client.Spawn()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") != "websocket" {
			http.FileServer(http.Dir(".")).ServeHTTP(w, r)
			return
		}
		websocket.Handler(func(ws *websocket.Conn) {
			ws.PayloadType = websocket.BinaryFrame
			conn := mux.New(ws)
			go func(src *mux.Session) {
				log.Fatal(mux.Proxy(c.Session, src))
			}(conn)
			conn.Wait()
		}).ServeHTTP(w, r)
	}))
	go http.ListenAndServe(":7778", nil)

	ctx := context.Background()
	_, err = c.Window.New(ctx, client.WindowOptions{
		Title:     "Electronish",
		Frameless: false,
		Visible:   true,
		Center:    true,
		URL:       "http://localhost:7778/cmd/electronish",
	})
	if err != nil {
		panic(err)
	}
	select {}
}
