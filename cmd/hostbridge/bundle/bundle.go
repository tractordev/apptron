package bundle

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"tractor.dev/hostbridge/bridge/misc"
	"tractor.dev/hostbridge/cmd/hostbridge/build"
	"tractor.dev/hostbridge/cmd/hostbridge/bundle/bundler"
	"tractor.dev/hostbridge/cmd/hostbridge/sign"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// currently specific to osx
func Bundle() {
	wd, err := os.Getwd()
	fatal(err)
	cmdname := filepath.Base(wd)
	appname := strings.Title(strings.Trim(cmdname, "_"))

	dir, err := ioutil.TempDir("", appname)
	fatal(err)
	defer os.RemoveAll(dir)
	//fmt.Println(dir)

	if _, err := os.Stat(filepath.Join(wd, cmdname)); os.IsNotExist(err) {
		build.Build()
	}
	os.MkdirAll(filepath.Join(dir, "assets"), 0755)
	exePath := filepath.Join(dir, "assets", appname)
	copyFile(filepath.Join(wd, cmdname), exePath)

	bundleName := fmt.Sprintf("com.progrium.%s", appname)
	err, errlog := sign.Sign(dir, bundleName, exePath)
	if err != nil {
		errlog.WriteTo(os.Stderr)
		log.Fatal(err)
	}

	bundler.AssetsDir = filepath.Dir(exePath)
	bundler.BinaryName = appname
	tmpfile, err := ioutil.TempFile("", "appicon.png")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	f, err := misc.Assets.Open("icon.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = io.Copy(tmpfile, f)
	if err != nil {
		log.Fatal(err)
	}
	bundler.IconFile = tmpfile.Name()
	// if _, err := os.Stat(filepath.Join(wd, "appicon.png")); os.IsNotExist(err) {
	// 	tmpfile, err := ioutil.TempFile("", "appicon.png")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer os.Remove(tmpfile.Name())
	// 	f, err := generate.Assets.Open("appicon.png")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer f.Close()
	// 	_, err = io.Copy(tmpfile, f)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	bundler.IconFile = tmpfile.Name()
	// } else {
	// 	bundler.IconFile = filepath.Join(wd, "appicon.png")
	// }
	bundler.BundleIdentifier = bundleName
	bundler.AppName = appname
	bundler.OutputDir = wd
	bundler.Bundle()
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
