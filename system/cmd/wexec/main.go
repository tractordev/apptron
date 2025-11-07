package main

import (
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
)

func debug(format string, a ...any) {
	if os.Getenv("DEBUG") == "1" {
		log.Printf(format+"\n", a...)
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	var taskType string
	flag.StringVar(&taskType, "type", "wasi", "Type of task")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("usage: wexec <wasm> [args...]")
	}

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
	// ultimately we shouldn't need to prefix the path with vm/1/fsys,
	// it should be relative to the task namespace
	args[0] = strings.TrimPrefix(filepath.Join("vm/1/fsys", wd, flag.Arg(0)), "/")

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
