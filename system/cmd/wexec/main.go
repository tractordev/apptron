package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	// "tractor.dev/wanix/fs/p9kit"
)

func debug(format string, a ...any) {
	if os.Getenv("DEBUG") == "1" {
		log.Printf(format+"\n", a...)
	}
}

func main() {
	// snuck in a couple debug commands to avoid a whole new executable

	// shmtest is a throughput test for the shared memory channel
	if len(os.Args) > 1 && os.Args[1] == "shmtest" {
		runShmTest()
		os.Exit(0)
	}

	// fuse is a test of the fuse filesystem
	if len(os.Args) > 1 && os.Args[1] == "fuse" {
		setupFuseFS()
		os.Exit(0)
	}

	// shm9p is 9p server of the root filesystem via shared memory
	if len(os.Args) > 1 && os.Args[1] == "shm9p" {
		runShm9P()
		os.Exit(0)
	}

	log.SetFlags(log.Lshortfile)
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("usage: wexec <wasm> [args...]")
	}

	taskType, err := detectWASMType(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	debug("detected WASM type: %s", taskType)

	// fake /env program to print environment for debugging
	if flag.Arg(0) == "/env" {
		fmt.Println(os.Environ())
		fmt.Println("---")
		for _, env := range os.Environ() {
			fmt.Println(">", env)
		}
		fmt.Println("---")
		fmt.Println(strings.Join(append(os.Environ(), ""), "\n"))
		os.Exit(0)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	args := flag.Args()
	absArg0, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	// ultimately we shouldn't need to prefix the path with vm/1/fsys,
	// it should be relative to the task namespace
	args[0] = strings.TrimPrefix(filepath.Join("vm/1/fsys", absArg0), "/")

	debug("allocating pid")
	pidRaw, err := os.ReadFile(fmt.Sprintf("/task/new/%s", taskType))
	if err != nil {
		log.Fatal(err)
	}
	pid := strings.TrimSpace(string(pidRaw))

	debug("writing cmd")
	if err := appendFile(fmt.Sprintf("/task/%s/cmd", pid), []byte(strings.Join(args, " "))); err != nil {
		log.Fatal(err)
	}

	debug("writing dir")
	if err := appendFile(fmt.Sprintf("/task/%s/dir", pid), []byte(wd)); err != nil {
		log.Fatal(err)
	}

	debug("writing env")
	env := strings.Join(append(os.Environ(), ""), "\n")
	if err := appendFile(fmt.Sprintf("/task/%s/env", pid), []byte(env)); err != nil {
		log.Fatal(err)
	}

	var done atomic.Int32
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		debug("polling fd/1 => stdout")
		f, err := os.Open(fmt.Sprintf("/task/%s/fd/1", pid))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		b := make([]byte, 4096)
		for {
			n, err := f.Read(b)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			if done.Load() == 1 && n == 0 {
				debug("stdout thread done")
				return
			}
			os.Stdout.Write(b[:n])
			time.Sleep(30 * time.Millisecond)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		debug("polling fd/2 => stderr")
		f, err := os.Open(fmt.Sprintf("/task/%s/fd/2", pid))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		b := make([]byte, 4096)
		for {
			n, err := f.Read(b)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			if done.Load() == 1 && n == 0 {
				debug("stderr thread done")
				return
			}
			os.Stderr.Write(b[:n])
			time.Sleep(30 * time.Millisecond)
		}
	}()

	debug("starting")
	if err := appendFile(fmt.Sprintf("/task/%s/ctl", pid), []byte("start")); err != nil {
		log.Fatal(err)
	}

	debug("waiting for exit")
	for {
		b, err := os.ReadFile(fmt.Sprintf("/task/%s/exit", pid))
		if err != nil {
			log.Fatal(err)
		}
		out := strings.TrimSpace(string(b))
		if out != "" {
			debug("exit code: %s", out)
			code, err := strconv.Atoi(out)
			if err != nil {
				log.Fatal(err)
			}
			done.Store(1)
			debug("waiting for threads to finish")
			wg.Wait()
			debug("exiting with code %d", code)
			os.Exit(code)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func appendFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func detectWASMType(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Skip WASM header (8 bytes: magic + version)
	f.Seek(8, 0)

	// Read sections until we find imports (section ID 2)
	for {
		var sectionID byte
		if err := binary.Read(f, binary.LittleEndian, &sectionID); err != nil {
			return "", err
		}

		size := readVarUint(f)

		if sectionID == 2 { // Import section
			buf := make([]byte, size)
			f.Read(buf)

			if bytes.Contains(buf, []byte("wasi_snapshot_preview1")) {
				return "wasi", nil
			}
			if bytes.Contains(buf, []byte("gojs")) {
				return "gojs", nil
			}
			return "", errors.New("unknown WASM type")
		}

		f.Seek(int64(size), io.SeekCurrent)
	}
}

func readVarUint(r io.Reader) uint64 {
	var v uint64
	var s uint
	b := []byte{0}
	for {
		r.Read(b)
		v |= uint64(b[0]&0x7f) << s
		if b[0]&0x80 == 0 {
			break
		}
		s += 7
	}
	return v
}
