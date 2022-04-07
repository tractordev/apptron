package misc

import (
	"context"
	"io/ioutil"

	"github.com/progrium/qtalk-go/rpc"
)

// FetchData is used in calls to fetch data blobs identified by a selector
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
