package apputil

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// TempCommand checks for the base name of cmdpath in PATH and if it doesn't find it,
// it will write out the file from fsys at cmdpath to a temporary directory as executable,
// returning the file path to it and a function to remove both.
//
// If the command is found in PATH or any step fails for any reason, TempCommand returns
// an empty string and a no-op function.
//
// This is intended to be used when embedding the hostbridge binary in a program, using this
// to write it out to disk temporarily, the path to which can be used to set the BRIDGECMD
// env var used by the hostbridge client Spawn function. It should also be cleaned up when
// the program finishes.
func TempCommand(fsys fs.FS, cmdpath string) (string, func()) {
	cmdname := filepath.Base(cmdpath)
	_, err := exec.LookPath(cmdname)
	if err == nil {
		return "", func() {}
	}
	f, err := fsys.Open(cmdpath)
	if err != nil {
		return "", func() {}
	}
	d, _ := ioutil.ReadAll(f)
	f.Close()
	dir, err := ioutil.TempDir("", fmt.Sprintf("%s-*", cmdname))
	if err != nil {
		return "", func() {}
	}
	path := filepath.Join(dir, cmdname)
	if err := ioutil.WriteFile(path, d, 0755); err != nil {
		return "", func() {}
	}
	return path, func() {
		os.RemoveAll(dir)
	}
}
