package client

import (
	"context"
	"fmt"
	"testing"
)

func testScreenModule(t *testing.T) {
	client, cleanup := setupBridgeClient(t)
	defer cleanup()

	ctx := context.Background()

	d, err := client.Screen.Displays(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(d)

}
