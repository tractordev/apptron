package client

import (
	"context"
	"io/ioutil"
	"path"
	"runtime"
	"testing"
)

func TestAppModule(t *testing.T) {
	client, cleanup := setupBridgeClient(t)
	defer cleanup()

	_, filename, _, _ := runtime.Caller(0)
	iconpath := path.Join(path.Dir(filename), "../assets/icon.png")
	icon, err := ioutil.ReadFile(iconpath)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	err = client.App.NewIndicator(ctx, icon, []MenuItem{
		{
			Title: "One",
		},
		{
			Title: "Two",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

}
