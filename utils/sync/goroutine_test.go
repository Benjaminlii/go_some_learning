package sync

import (
	"context"
	"fmt"
	"testing"
)

func TestDoWithGoroutine(t *testing.T) {
	req := []interface{}{}
	for i := 0; i < 10000; i++ {
		req = append(req, i)
	}
	resp, err := DoWithGoroutine(context.Background(),
		func(ctx context.Context, reqList []interface{}) DoWithGoroutineResp {
			fmt.Printf("       gorouting print : %v\n", reqList)
			return DoWithGoroutineResp{
				Resp: reqList,
				Err:  nil,
			}
		},
		req,
		3,
	)

	if err != nil {
		fmt.Println("err")
	} else {
		fmt.Println(len(resp))
	}
}
