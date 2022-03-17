package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	_ "embed"
)

//go:embed app.go.res
var entrypoint []byte

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	os.Setenv("GOPRIVATE", "tractor.dev/*")
	gobin, err := exec.LookPath("go")
	fatal(err)

	dir, err := os.Getwd()
	fatal(err)
	appname := filepath.Base(dir)

	if _, err := os.Stat("go.mod"); err == nil {
		log.Fatal("unable to run in dir with go.mod")
	}

	if _, err := os.Stat("main.go"); err == nil {
		log.Fatal("unable to run in dir with main.go")
	}

	fatal(ioutil.WriteFile("main.go", entrypoint, 0644))

	run(gobin, "mod", "init", appname)
	run(gobin, "get", "-u", "tractor.dev/hostbridge")
	run(gobin, "get")
	run(gobin, "build", "-o", appname, ".")

	fatal(os.Remove("main.go"))
	fatal(os.Remove("go.mod"))
	fatal(os.Remove("go.sum"))

}

func run(args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
