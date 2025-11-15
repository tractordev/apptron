package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"tractor.dev/wanix/vm/v86/shm"
)

func runShmTest() {
	sch, err := shm.NewSharedChannel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create channel: %v\n", err)
		os.Exit(1)
	}
	defer sch.Close()

	fmt.Println("Shared channel connected!")

	fmt.Println("Starting throughput test...")
	if err := ThroughputTest(sch); err != nil {
		fmt.Fprintf(os.Stderr, "Throughput test failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Throughput test completed successfully")

}

// ThroughputTest measures roundtrip (echo) throughput via the provided stream.
// It writes data to the stream and reads the corresponding echoed data back,
// reporting the transfer rate and latency.
func ThroughputTest(conn io.ReadWriteCloser) error {
	const (
		testSize  = 1024 * 1024 * 80 // 80 MB total transfer size
		chunkSize = 256 * 1024       // 256KB per write
	)
	sendBuf := make([]byte, chunkSize)
	for i := range sendBuf {
		sendBuf[i] = byte(i)
	}
	recvBuf := make([]byte, chunkSize)
	// nwrites := testSize / chunkSize

	start := time.Now()
	totalSent := 0
	totalRecv := 0

	go func() {
		// Write phase. This is actually asynchronous.
		for totalSent < testSize {
			n, err := conn.Write(sendBuf)
			if err != nil {
				log.Fatalf("write error: %v", err)
			}
			totalSent += n
		}
	}()

	// Read phase (echo expected)
	i := 0
	for totalRecv < testSize {
		n, err := conn.Read(recvBuf)
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}
		totalRecv += n
		if i%100 == 0 {
			fmt.Println("inner: read:", n, "totalRecv:", totalRecv, "totalSent:", totalSent, "i:", i)
		}
		i++
	}

	elapsed := time.Since(start)
	mb := float64(totalSent) / (1024 * 1024)
	mbps := mb / elapsed.Seconds()
	fmt.Printf("Throughput: sent %d bytes, recv %d bytes in %.3fs (%.2f MB/s)\n",
		totalSent, totalRecv, elapsed.Seconds(), mbps)

	return nil
}
