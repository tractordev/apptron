package client

import (
	"context"
	"fmt"
	"testing"
)

func TestMenuModule(t *testing.T) {
	client, cleanup := setupBridgeClient(t)
	defer cleanup()

	ctx := context.Background()

	m, err := client.Menu.New(ctx, []MenuItem{
		{
			Title:   "One",
			Enabled: true,
		},
		{
			Title:   "Two",
			Enabled: true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m)

}
