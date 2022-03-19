//go:build !entrypoint

package demo

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	_ "embed"
)

//go:embed entry.go
var entrypoint []byte

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Build() {
	os.Setenv("GOPRIVATE", "tractor.dev/*")
	gobin, err := exec.LookPath("go")
	fatal(err)
	selfbin, err := os.Executable()
	fatal(err)

	dir, err := os.Getwd()
	fatal(err)
	appname := filepath.Base(dir)

	for _, name := range []string{"go.mod", "main.go", "hostbridge"} {
		if _, err := os.Stat(name); err == nil {
			log.Fatalf("unable to run in dir with %s", name)
		}
	}

	fatal(copyFile(selfbin, filepath.Join(dir, "hostbridge")))
	fatal(ioutil.WriteFile("main.go", entrypoint, 0644))

	run(gobin, "mod", "init", appname)
	run(gobin, "get", "-u", "tractor.dev/hostbridge")
	run(gobin, "get")
	run(gobin, "build", "-tags", "entrypoint", "-o", appname, ".")

	for _, name := range []string{"go.mod", "go.sum", "main.go", "hostbridge"} {
		fatal(os.Remove(name))
	}

}

func run(args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func copyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
