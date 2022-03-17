package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	gobin, err := exec.LookPath("go")
	if err != nil {
		log.Fatal(err)
	}

	// write main.go into curr dir

	cmd := exec.Command(gobin, "build", "-o", "app", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

}
