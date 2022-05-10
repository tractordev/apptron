package cli

import (
	"context"
	"io"
)

// ContextWithIO returns a child context with a ContextIO
// value added using the given Stdio equivalents.
func ContextWithIO(parent context.Context, in io.Reader, out io.Writer, err io.Writer) context.Context {
	return context.WithValue(parent, IOContextKey, &iocontext{
		in:  in,
		out: out,
		err: err,
	})
}

type iocontext struct {
	out, err io.Writer
	in       io.Reader
}

func (c *iocontext) Write(p []byte) (n int, err error) {
	return c.out.Write(p)
}

func (c *iocontext) Read(p []byte) (n int, err error) {
	return c.in.Read(p)
}

func (c *iocontext) Err() io.Writer {
	return c.err
}

// ContextIO is an io.ReadWriter with an extra io.Writer
// for an error channel. Typically wrapping STDIO.
type ContextIO interface {
	io.Reader
	io.Writer
	Err() io.Writer
}

// IOContextKey is the context key used by IOFrom to get a ContextIO.
var IOContextKey = "io"

// IOFrom pulls a ContextIO from a context.
func IOFrom(ctx context.Context) ContextIO {
	v := ctx.Value(IOContextKey)
	if v == nil {
		return nil
	}
	return v.(ContextIO)
}
