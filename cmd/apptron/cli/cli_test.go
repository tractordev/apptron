package cli

import (
	"bytes"
	"context"
	"fmt"
	"testing"
)

// Tests TODO:
// top level help
// command help
// hidden commands
// command aliases
// help, examples

func TestSimpleCommand(t *testing.T) {
	cmd := &Command{
		Usage: "simple",
		Run: func(ctx context.Context, args []string) {
			fmt.Fprint(IOFrom(ctx), "Hello world")
		},
	}
	var buf bytes.Buffer
	ctx := ContextWithIO(context.Background(), nil, &buf, nil)
	if err := Execute(ctx, cmd, []string{}); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "Hello world" {
		t.Fatal("unexpected output")
	}
}

func TestPositionalArgs(t *testing.T) {
	cmd := &Command{
		Usage: "posargs",
		Args:  ExactArgs(2),
		Run:   func(ctx context.Context, args []string) {},
	}
	if err := Execute(context.Background(), cmd, []string{"one", "two"}); err != nil {
		t.Fatal(err)
	}
	if err := Execute(context.Background(), cmd, []string{"one", "two", "three"}); err == nil {
		t.Fatal("expected error")
	}
	if err := Execute(context.Background(), cmd, []string{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestSubcommands(t *testing.T) {
	cmd := &Command{
		Usage: "subcmds",
		Run: func(ctx context.Context, args []string) {
			fmt.Fprint(IOFrom(ctx), "root")
		},
	}
	cmd.AddCommand(&Command{
		Usage: "sub1",
		Run: func(ctx context.Context, args []string) {
			fmt.Fprint(IOFrom(ctx), "sub1")
		},
	})
	cmd.AddCommand(&Command{
		Usage: "sub2",
		Run: func(ctx context.Context, args []string) {
			fmt.Fprint(IOFrom(ctx), "sub2")
		},
	})

	var buf bytes.Buffer
	ctx := ContextWithIO(context.Background(), nil, &buf, nil)

	if err := Execute(ctx, cmd, []string{}); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "root" {
		t.Fatal("didn't run root cmd")
	}

	buf.Reset()
	if err := Execute(ctx, cmd, []string{"sub1"}); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "sub1" {
		t.Fatal("didn't run sub1 cmd")
	}

	buf.Reset()
	if err := Execute(ctx, cmd, []string{"sub2"}); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "sub2" {
		t.Fatal("didn't run sub2 cmd")
	}
}

func TestFlags(t *testing.T) {
	var (
		boolFlag   bool
		stringFlag string
		intFlag    int
	)
	cmd := &Command{
		Usage: "flags",
		Run:   func(ctx context.Context, args []string) {},
	}
	cmd.Flags().BoolVar(&boolFlag, "b", false, "bool value")
	cmd.Flags().StringVar(&stringFlag, "string", "", "string value")
	cmd.Flags().IntVar(&intFlag, "int", 0, "int value")

	ctx := context.Background()
	if err := Execute(ctx, cmd, []string{"-b=true", "-string", "STRING", "-int=100"}); err != nil {
		t.Fatal(err)
	}
	if boolFlag != true || stringFlag != "STRING" || intFlag != 100 {
		t.Fatal("unexpected flag value")
	}
}

func TestVersion(t *testing.T) {
	cmd := &Command{
		Version: "1.0",
		Usage:   "mytest",
		Run:     func(ctx context.Context, args []string) {},
	}
	var buf bytes.Buffer
	ctx := ContextWithIO(context.Background(), nil, &buf, nil)
	if err := Execute(ctx, cmd, []string{"-v"}); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "1.0\n" {
		t.Fatal("unexpected output:", buf.String())
	}
}
