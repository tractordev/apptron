package core

import "C"
import (
	"context"
	"io/ioutil"

	"github.com/progrium/qtalk-go/rpc"
)

type Position struct {
	X float64
	Y float64
}

type Size struct {
	Width  float64
	Height float64
}

type Handle int

func FetchData(ctx context.Context, call *rpc.Call, selector string) ([]byte, error) {
	resp, err := call.Caller.Call(ctx, selector, nil, nil)
	if err != nil {
		return nil, err
	}
	ch := resp.Channel
	defer ch.Close()
	data, err := ioutil.ReadAll(ch)
	if err != nil {
		return nil, err
	}
	return data, nil
}
