package main

import (
	"context"
	"embed"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/terminfo"
	"github.com/progrium/qtalk-go/mux"
	"github.com/progrium/qtalk-go/rpc"
	"github.com/progrium/qtalk-go/x/cbor/codec"
	"github.com/rivo/tview"
	"golang.org/x/net/websocket"
	"tractor.dev/apptron"
)

//go:embed *.html xterm qtalk
var assets embed.FS

func main() {
	ctx := context.Background()

	native, err := apptron.Run(ctx, apptron.AppOptions{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = native.Window.New(ctx, apptron.WindowOptions{
		Visible:   true,
		Size:      apptron.Size{Width: 800, Height: 600},
		Center:    true,
		Resizable: true,
		Title:     "TUI Demo",
		URL:       "http://localhost:9090",
	})
	if err != nil {
		log.Fatal(err)
	}

	var term *tcellTTY

	methods := rpc.NewRespondMux()
	methods.Handle("terminal", rpc.HandlerFunc(func(r rpc.Responder, c *rpc.Call) {
		var pi ptyInfo
		if err := c.Receive(&pi); err != nil {
			panic(err)
		}
		ch, err := r.Continue(nil)
		if err != nil {
			panic(err)
		}

		term = &tcellTTY{
			ReadWriteCloser: ch,
			info:            pi,
		}
		ti, err := terminfo.LookupTerminfo(pi.Term)
		if err != nil {
			panic(err)
		}
		screen, err := tcell.NewTerminfoScreenFromTtyTerminfo(term, ti)
		if err != nil {
			panic(err)
		}

		// tview says we don't have to do this
		// when using SetScreen, but it lies
		if err := screen.Init(); err != nil {
			panic(err)
		}

		app := tview.NewApplication().SetScreen(screen).EnableMouse(true)

		modal := tview.NewModal().
			SetText("Do you want to quit the application?").
			AddButtons([]string{"Quit", "Cancel"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Quit" {
					app.Stop()
					native.Close()
				}
			})

		app.SetRoot(modal, false)
		if err := app.Run(); err != nil {
			panic(err)
		}

		ch.Close()
	}))
	methods.Handle("terminal.resize", rpc.HandlerFunc(func(r rpc.Responder, c *rpc.Call) {
		c.Receive(&term.info)
		term.notifyResize()
	}))

	http.Handle("/", http.FileServer(http.FS(assets)))
	http.Handle("/rpc", websocket.Handler(func(ws *websocket.Conn) {
		ws.PayloadType = websocket.BinaryFrame
		conn := mux.New(ws)
		srv := &rpc.Server{
			Handler: methods,
			Codec:   codec.CBORCodec{},
		}
		go srv.Respond(conn, context.Background())
		conn.Wait()
	}))
	go http.ListenAndServe(":9090", nil)

	native.Wait()
}

type ptyInfo struct {
	Term string
	Cols int
	Rows int
}

type tcellTTY struct {
	io.ReadWriteCloser
	resizecb func()
	mu       sync.Mutex
	info     ptyInfo
}

func (t *tcellTTY) Start() error {
	return nil
}

func (t *tcellTTY) Stop() error {
	return nil
}

func (t *tcellTTY) Drain() error {
	return nil
}

func (t *tcellTTY) WindowSize() (int, int, error) {
	return t.info.Cols, t.info.Rows, nil
}

func (t *tcellTTY) NotifyResize(cb func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.resizecb = cb
}

func (t *tcellTTY) notifyResize() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.resizecb != nil {
		t.resizecb()
	}
}
