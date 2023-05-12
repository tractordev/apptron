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
		// @Incomplete:
		// - menu.Menu needs to be replaced with "Menu"
		//{"bridge/api/app/app.go", "client/app.go", "App"},

		// @Incomplete:
		// - lowercase "menu" needs to be replaced with uppercase "Menu"
		// - resource.Handle gets replaced with Handle
		//{"bridge/api/menu/menu.go", "client/menu.go", "Menu"},

		{"bridge/api/shell/shell.go", "client/shell.go", "Shell"},
		{"bridge/api/system/system.go", "client/system.go", "System"},

		//{"bridge/api/window/window.go", "client/window.go", "Window"},
	}

	filePrefix := `package client

import (
	"context"
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
		log.Println(result)
	}
}
