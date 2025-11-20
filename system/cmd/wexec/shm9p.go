package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hugelgupf/p9/p9"
	"tractor.dev/wanix/fs/localfs"
	"tractor.dev/wanix/fs/p9kit"
	"tractor.dev/wanix/vm/v86/shm"
)

func runShm9P() {
	sch, err := shm.NewSharedChannel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create channel: %v\n", err)
		os.Exit(1)
	}
	defer sch.Close()

	go monitorNewPorts(1 * time.Second)

	dirfs, err := localfs.New("/")
	if err != nil {
		log.Fatal(err)
	}
	p9srv := p9.NewServer(p9kit.Attacher(dirfs)) //, p9.WithServerLogger(ulog.Log))
	go func() {
		err := os.WriteFile("/run/shm9p.lock", []byte(""), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}()
	defer func() {
		os.Remove("/run/shm9p.lock")
	}()
	if err := p9srv.Handle(sch, sch); err != nil {
		log.Fatal(err)
	}
}

type ListeningPort struct {
	Port    int
	Address string
}

func getListeningPorts() (map[int]bool, error) {
	ports := make(map[int]bool)

	// Check both IPv4 and IPv6
	for _, file := range []string{"/proc/net/tcp"} {
		if err := parseNetFile(file, ports); err != nil {
			return nil, err
		}
	}

	return ports, nil
}

func parseNetFile(filename string, ports map[int]bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // Skip header

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}

		// State 0A = LISTEN
		if fields[3] == "0A" {
			// Parse local address
			parts := strings.Split(fields[1], ":")
			if len(parts) == 2 {
				// Convert hex port to decimal
				if port, err := strconv.ParseInt(parts[1], 16, 32); err == nil {
					ports[int(port)] = true
				}
			}
		}
	}

	return scanner.Err()
}

func monitorNewPorts(interval time.Duration) {
	knownPorts, _ := getListeningPorts()
	fmt.Printf("Initial listening ports: %v\n", knownPorts)

	for {
		time.Sleep(interval)

		currentPorts, err := getListeningPorts()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Check for new ports
		for port := range currentPorts {
			if !knownPorts[port] {
				fmt.Printf("🔴 NEW PORT LISTENING: %d\n", port)
				knownPorts[port] = true
			}
		}

		// Check for closed ports
		for port := range knownPorts {
			if !currentPorts[port] {
				fmt.Printf("⚪ PORT CLOSED: %d\n", port)
				delete(knownPorts, port)
			}
		}
	}
}
