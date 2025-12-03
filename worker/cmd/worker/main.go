package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
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
		// handle bundles
		if strings.HasPrefix(r.URL.Path, "/bundles/") {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			if strings.HasSuffix(r.URL.Path, ".gz") {
				w.Header().Set("Content-Encoding", "gzip")
			} else {
				w.Header().Set("Content-Encoding", "br")
			}
			http.StripPrefix("/bundles/", http.FileServer(http.Dir("bundles"))).ServeHTTP(w, r)
			return
		}

		// rest are network requests, so make sure the network is available
		if vn == nil {
			http.Error(w, "network not available", http.StatusNotFound)
			return
		}

		// handle port tunnel requests
		porthost := r.Host
		if r.URL.Query().Has("port") {
			porthost = r.URL.Query().Get("port")
			q := r.URL.Query()
			q.Del("port")
			r.URL.RawQuery = q.Encode()
		}
		if strings.HasPrefix(porthost, "tcp-") {
			parts := strings.Split(porthost, ".")
			port := strings.TrimPrefix(parts[0], "tcp-")
			ip := parts[1]
			if strings.Contains(ip, "-") {
				ip = strings.Replace(ip, "-", ".", -1)
			} else {
				var err error
				ip, err = DecodeIP(ip)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

			conn, err := vn.Dial("tcp", net.JoinHostPort(ip, port))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer conn.Close()

			u, err := url.Parse(fmt.Sprintf("http://%s", net.JoinHostPort(ip, port)))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			proxy := CreateProxyWithConn(u, conn)
			proxy.ServeHTTP(w, r)
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

// DecodeIP converts "HHHHHHHH" hex to "IP"
func DecodeIP(encoded string) (string, error) {
	ipBytes, err := hex.DecodeString(encoded)
	if err != nil || len(ipBytes) != 4 {
		return "", fmt.Errorf("invalid IP hex")
	}
	ip := net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
	return ip.String(), nil
}

// EncodeIP converts an IPv4 string (e.g. "127.0.0.1") to its "HHHHHHHH" hex format.
func EncodeIP(ipstr string) (string, error) {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return "", fmt.Errorf("invalid IP address")
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return "", fmt.Errorf("not an IPv4 address")
	}
	return hex.EncodeToString(ipv4), nil
}

// CustomDialer wraps a specific net.Conn to be used by the HTTP transport
type CustomDialer struct {
	conn net.Conn
}

func (d *CustomDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	// Return the pre-established connection
	// Note: This simple implementation returns the same conn every time
	// You may want to handle connection reuse/pooling differently
	return d.conn, nil
}

// CreateProxyWithConn creates a reverse proxy that uses a specific net.Conn
func CreateProxyWithConn(targetURL *url.URL, conn net.Conn) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Create a custom transport that uses our specific connection
	transport := &http.Transport{
		DialContext: (&CustomDialer{conn: conn}).DialContext,
		// Disable connection pooling since we're managing the connection manually
		MaxIdleConns:        1,
		MaxIdleConnsPerHost: 1,
		DisableKeepAlives:   false,
		IdleConnTimeout:     90 * time.Second,
	}

	proxy.Transport = transport

	return proxy
}
