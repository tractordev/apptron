package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"tractor.dev/toolkit-go/engine/cli"
)

func portsCmd() *cli.Command {
	return &cli.Command{
		Usage: "ports",
		Short: "monitor listening ports",
		Run:   monitorPorts,
	}
}

func monitorPorts(ctx *cli.Context, args []string) {
	interval := time.Duration(1 * time.Second)
	knownPorts, _ := getListeningPorts()

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
				url := portURL(port)
				if url != "" {
					fmt.Printf("\n=> Apptron public URL: %s\n\n", url)
				}
				knownPorts[port] = true
			}
		}

		// Check for closed ports
		for port := range knownPorts {
			if !currentPorts[port] {
				// fmt.Printf("=> Apptron port closed: %d\n", port)
				delete(knownPorts, port)
			}
		}
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

// encodeIP converts an IPv4 string (e.g. "127.0.0.1") to its "HHHHHHHH" hex format.
func encodeIP(ipstr string) (string, error) {
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

func portURL(port int) string {
	sessionIP := os.Getenv("SESSION_IP")
	if sessionIP == "" {
		return ""
	}
	user := os.Getenv("USER")
	if user == "" {
		return ""
	}
	ip, err := encodeIP(sessionIP)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("https://tcp-%d-%s-%s.apptron.dev", port, ip, user)
}
