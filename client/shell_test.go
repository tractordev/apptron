package client

import (
	"context"
	"fmt"
	"testing"
)

func TestShellModule(t *testing.T) {
	client, cleanup := setupBridgeClient(t)
	defer cleanup()

	ctx := context.Background()

	err := client.Shell.ShowNotification(ctx, Notification{
		Title:    "Hello Title",
		Subtitle: "Hello Subtitle",
		Body:     "Hello body",
	})
	if err != nil {
		t.Fatal(err)
	}

	ret, err := client.Shell.ShowMessage(ctx, MessageDialog{
		Title:   "Test Title",
		Body:    "Test Body",
		Buttons: "ok",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ret)

	folder, err := client.Shell.ShowFilePicker(ctx, FileDialog{
		Title:     "Pick a file",
		Directory: "/Users",
		Mode:      "pickfolder",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(folder)

	text := "Hello clipboard"
	_, err = client.Shell.WriteClipboard(ctx, text)
	if err != nil {
		t.Fatal(err)
	}
	cliptext, err := client.Shell.ReadClipboard(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if cliptext != text {
		t.Fatalf("unexpected clipboard value: %s", cliptext)
	}

	client.Shell.OnShortcut = func(e Event) {
		fmt.Println(e)
	}
	err = client.Shell.RegisterShortcut(ctx, "Control+Shift+T")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Shell.UnregisterAllShortcuts(ctx)
	if err != nil {
		t.Fatal(err)
	}

}
