package build

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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

	parts := strings.Split(dir, string(filepath.Separator))
	workdir := filepath.Join(os.TempDir(), strings.Join(parts[len(parts)-3:], "-"))
	os.MkdirAll(workdir, 0755)

	fmt.Printf("building %s ...\n", appname)

	start := time.Now()

	binFile := filepath.Join(workdir, "hostbridge")
	if _, err := os.Stat(binFile); err != nil {
		fatal(copyFile(selfbin, binFile))
	}

	mainFile := filepath.Join(workdir, "main.go")
	if _, err := os.Stat(mainFile); err != nil {
		fatal(ioutil.WriteFile(mainFile, entrypoint[19:], 0644))
	}

	di, err := ioutil.ReadDir(dir)
	fatal(err)
	for _, fi := range di {
		if !fi.IsDir() && fi.Name() != appname {
			fatal(copyFile(filepath.Join(dir, fi.Name()), filepath.Join(workdir, fi.Name())))
		}
	}

	var buf bytes.Buffer

	modFile := filepath.Join(workdir, "go.mod")
	if _, err := os.Stat(modFile); err != nil {
		run(&buf, workdir, gobin, "mod", "init", appname)
		run(&buf, workdir, gobin, "get", "-u", "tractor.dev/hostbridge")
		run(&buf, workdir, gobin, "get")
	}

	run(&buf, workdir, gobin, "build", "-o", filepath.Join(dir, appname), ".")

	fmt.Printf("done! [%s]\n", time.Since(start).Round(time.Millisecond))

}

func run(buf *bytes.Buffer, dir string, args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = buf
	cmd.Stderr = buf
	if err := cmd.Run(); err != nil {
		buf.WriteTo(os.Stderr)
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
