package main

import (
	"context"
	"fmt"

	"github.com/Benjaminlii/go_some_learning/utils/sync"
)

func main() {
	req := []interface{}{}
	for i := 0; i < 10000; i++ {
		req = append(req, i)
	}
	resp, err := sync.DoWithGoroutine(context.Background(),
		func(ctx context.Context, reqList []interface{}) sync.DoWithGoroutineResp {
			fmt.Printf("       gorouting print : %v\n", reqList)
			return sync.DoWithGoroutineResp{
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
