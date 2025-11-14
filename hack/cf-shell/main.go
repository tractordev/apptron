package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/gorilla/websocket"
	"tractor.dev/wanix/fs/fskit"
	"tractor.dev/wanix/fs/fusekit"
)

// Linux-specific ioctl constants for PTY
const (
	TIOCGPTN   = 0x80045430
	TIOCSPTLCK = 0x40045431
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Spawn shell (try bash first, fall back to sh)
	shell := "/bin/bash"
	if _, err := os.Stat(shell); os.IsNotExist(err) {
		shell = "/bin/sh"
	}
	cmd := exec.Command(shell)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	// Try to create a pseudo-terminal
	ptmx, err := openPty(cmd)
	if err != nil {
		// Fall back to non-PTY mode
		handleWebSocketNoPty(conn, shell)
		return
	}
	defer ptmx.Close()

	// Start the shell
	if err := cmd.Start(); err != nil {
		ptmx.Close()
		// Fall back to non-PTY mode
		handleWebSocketNoPty(conn, shell)
		return
	}

	// Handle cleanup
	done := make(chan struct{})
	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
		close(done)
	}()

	// Read from pty and write to websocket
	go func() {
		buf := make([]byte, 8192)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	// Read from websocket and write to pty
	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if msgType == websocket.TextMessage {
			if _, err := ptmx.Write(data); err != nil {
				break
			}
		}
	}
}

func handleWebSocketNoPty(conn *websocket.Conn, shell string) {
	// Fallback method using pipes instead of PTY
	cmd := exec.Command(shell)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color", "PS1=\\$ ", "HOME=/root")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	if err := cmd.Start(); err != nil {
		return
	}

	// Send initial prompt since we don't have PTY
	conn.WriteMessage(websocket.TextMessage, []byte("$ "))

	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	// Read from stdout and write to websocket
	go func() {
		buf := make([]byte, 8192)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	// Read from stderr and write to websocket
	go func() {
		buf := make([]byte, 8192)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	// Read from websocket and write to stdin
	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if msgType == websocket.TextMessage {
			if _, err := stdin.Write(data); err != nil {
				break
			}
		}
	}
}

func openPty(cmd *exec.Cmd) (*os.File, error) {
	// Open a new pseudo-terminal
	ptmx, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open /dev/ptmx: %w", err)
	}

	// Grant access to the slave pty
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, ptmx.Fd(), TIOCGPTN, uintptr(unsafe.Pointer(new(int)))); errno != 0 {
		ptmx.Close()
		return nil, fmt.Errorf("TIOCGPTN ioctl failed: %v", errno)
	}

	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, ptmx.Fd(), TIOCSPTLCK, uintptr(unsafe.Pointer(new(int)))); errno != 0 {
		ptmx.Close()
		return nil, fmt.Errorf("TIOCSPTLCK ioctl failed: %v", errno)
	}

	// Get the name of the slave pty
	ptsname, err := ptsname(ptmx)
	if err != nil {
		ptmx.Close()
		return nil, err
	}

	// Open the slave pty
	pts, err := os.OpenFile(ptsname, os.O_RDWR, 0)
	if err != nil {
		ptmx.Close()
		return nil, fmt.Errorf("failed to open slave pty: %w", err)
	}

	// Set up the command to use the slave pty
	cmd.Stdin = pts
	cmd.Stdout = pts
	cmd.Stderr = pts
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid:  true,
		Setctty: true,
	}

	// Close pts in parent after cmd starts
	go func() {
		cmd.Wait()
		pts.Close()
	}()

	return ptmx, nil
}

func ptsname(f *os.File) (string, error) {
	var n int
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), TIOCGPTN, uintptr(unsafe.Pointer(&n)))
	if errno != 0 {
		return "", errno
	}
	return fmt.Sprintf("/dev/pts/%d", n), nil
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Web Terminal</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/xterm@5.3.0/css/xterm.css" />
    <style>
        body {
            margin: 0;
            padding: 0;
            background-color: #000;
            overflow: hidden;
        }
        #terminal {
            width: 100vw;
            height: 100vh;
            padding: 10px;
            box-sizing: border-box;
        }
    </style>
</head>
<body>
    <div id="terminal"></div>
    <script src="https://cdn.jsdelivr.net/npm/xterm@5.3.0/lib/xterm.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/xterm-addon-fit@0.8.0/lib/xterm-addon-fit.js"></script>
    <script>
        // Create terminal
        const term = new Terminal({
            cursorBlink: true,
            fontSize: 14,
            fontFamily: 'Menlo, Monaco, "Courier New", monospace',
            theme: {
                background: '#000000',
                foreground: '#ffffff'
            }
        });

        // Add fit addon for fullscreen
        const fitAddon = new FitAddon.FitAddon();
        term.loadAddon(fitAddon);

        // Open terminal in the container
        term.open(document.getElementById('terminal'));
        fitAddon.fit();

        // Handle window resize
        window.addEventListener('resize', () => {
            fitAddon.fit();
        });

        // Connect to WebSocket
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const ws = new WebSocket(protocol + '//' + window.location.host + '/ws');

        ws.onopen = () => {
            console.log('WebSocket connected');
            
            // Send data from terminal to websocket
            term.onData(data => {
                ws.send(data);
            });
        };

        ws.onmessage = (event) => {
            // Write data from websocket to terminal
            term.write(event.data);
        };

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            term.write('\r\n\x1b[31mWebSocket error occurred\x1b[0m\r\n');
        };

        ws.onclose = () => {
            console.log('WebSocket closed');
            term.write('\r\n\x1b[31mConnection closed\x1b[0m\r\n');
        };

        // Focus terminal on load
        term.focus();
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func main() {
	fsys := fskit.MapFS{
		"helloworld.txt": fskit.RawNode([]byte("Hello, World!")),
	}
	mount, err := fusekit.Mount(fsys, "/x", context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mount.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/ws", handleWebSocket)

	addr := ":8080"
	log.Printf("Starting web terminal server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
