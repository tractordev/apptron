package main

import (
	"fmt"
	"log"
	"os"
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
