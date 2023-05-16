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

/*

func Bind(client *Client, name string, f reflect.Value) {
	f.Set(reflect.MakeFunc(f.Type(), func(args []reflect.Value) (result []reflect.Value) {
		ctx, _ := reflect.ValueOf(args[0]).Interface().(context.Context)
		_, err := client.Call(ctx, name, fn.Args{args[1]}, nil)
		return []reflect.Value{reflect.ValueOf(&err).Elem()}
	}))
}

func BindAll(client *Client, name string, p interface{}) {
	value := reflect.ValueOf(p).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)

		if field.Type.Kind() == reflect.Func && !strings.HasPrefix(field.Name, "On") {
			fullName := name + "." + field.Name
			log.Println(field.Index[0], field.Name, fullName)
			Bind(client, fullName, value.Field(i))
		}
	}
}

*/

func main() {
	gens := []gen{
		// @Incomplete:
		// - menu.Menu needs to be replaced with "Menu"
		//{"bridge/api/app/app.go", "client/app.go", "App"},

		// @Incomplete:
		// - lowercase "menu" needs to be replaced with uppercase "Menu"
		// - resource.Handle gets replaced with Handle
		//{"bridge/api/menu/menu.go", "client/menu.go", "Menu"},

		{"bridge/api/shell/shell.go", "schema/out/shell.go", "Shell"},
		{"bridge/api/system/system.go", "schema/out/system.go", "System"},

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
