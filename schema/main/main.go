package main

import (
	"log"
	"os"

	"tractor.dev/apptron/schema"
)

type gen struct {
	InputPath  string
	OutputPath string
	ModName    string
}

func main() {
	/*
		files := []string{
			"bridge/api/app/app.go",
		}
	*/

	gens := []gen{
		//{"bridge/api/app/app.go", "client/app.go", "App"},
		{"bridge/api/system/system.go", "client/system.go", "System"},
	}

	filePrefix := `package client

import (
	"context"

	"github.com/progrium/qtalk-go/fn"
)
`

	for _, gen := range gens {
		result := schema.GenerateClientCode(gen.InputPath, gen.ModName)

		f, err := os.Create(gen.OutputPath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		_, err = f.WriteString(filePrefix + "\n" + result)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Wrote", gen.OutputPath)
	}
}
