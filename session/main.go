package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/progrium/go-netstack/vnet"
)

func main() {
	vn, err := vnet.New(&vnet.Configuration{
		Debug:             false,
		MTU:               1500,
		Subnet:            "192.168.127.0/24",
		GatewayIP:         "192.168.127.1",
		GatewayMacAddress: "5a:94:ef:e4:0c:dd",
		GatewayVirtualIPs: []string{"192.168.127.253"},
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8080", handler(vn)); err != nil {
		log.Fatal(err)
	}
}

func handler(vn *vnet.VirtualNetwork) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if vn == nil {
			http.Error(w, "network not available", http.StatusNotFound)
			return
		}
		if !websocket.IsWebSocketUpgrade(r) {
			http.Error(w, "expecting websocket upgrade", http.StatusBadRequest)
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer ws.Close()

		fmt.Println("network session started")

		if err := vn.AcceptQemu(r.Context(), &qemuAdapter{Conn: ws}); err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				return
			}
			log.Println(err)
			return
		}
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type qemuAdapter struct {
	*websocket.Conn
	mu          sync.Mutex
	readBuffer  []byte
	writeBuffer []byte
	readOffset  int
}

func (q *qemuAdapter) Read(p []byte) (n int, err error) {
	if len(q.readBuffer) == 0 {
		_, message, err := q.ReadMessage()
		if err != nil {
			return 0, err
		}
		length := uint32(len(message))
		lengthPrefix := make([]byte, 4)
		binary.BigEndian.PutUint32(lengthPrefix, length)
		q.readBuffer = append(lengthPrefix, message...)
		q.readOffset = 0
	}

	n = copy(p, q.readBuffer[q.readOffset:])
	q.readOffset += n
	if q.readOffset >= len(q.readBuffer) {
		q.readBuffer = nil
	}
	return n, nil
}

func (q *qemuAdapter) Write(p []byte) (int, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.writeBuffer = append(q.writeBuffer, p...)

	if len(q.writeBuffer) < 4 {
		return len(p), nil
	}

	length := binary.BigEndian.Uint32(q.writeBuffer[:4])
	if len(q.writeBuffer) < int(length)+4 {
		return len(p), nil
	}

	err := q.WriteMessage(websocket.BinaryMessage, q.writeBuffer[4:4+length])
	if err != nil {
		return 0, err
	}

	q.writeBuffer = q.writeBuffer[4+length:]
	return len(p), nil
}

func (c *qemuAdapter) LocalAddr() net.Addr {
	return &net.UnixAddr{}
}

func (c *qemuAdapter) RemoteAddr() net.Addr {
	return &net.UnixAddr{}
}

func (c *qemuAdapter) SetDeadline(t time.Time) error {
	return nil
}
func (c *qemuAdapter) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *qemuAdapter) SetWriteDeadline(t time.Time) error {
	return nil
}
