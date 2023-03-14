package sync

import (
	"context"
	"math"
	"runtime/debug"
	"sync/atomic"

	"github.com/pkg/errors"

	"github.com/Benjaminlii/go_some_learning/common/logger"
)

func recoverAndLogWithCtx(ctx context.Context) {
	if err := recover(); err != nil {
		stack := string(debug.Stack())
		logger.Errorf(ctx, "[SafeGo] Goroutine Recover: %+v, stack is %v", err, stack)
	}
}

func safeGoWithCtx(ctx context.Context, f func()) {
	defer recoverAndLogWithCtx(ctx)
	f()
}

func worker(ctx context.Context, f WorkerFunc, in <-chan interface{}, goroutineNum int32) <-chan interface{} {
	out := make(chan interface{}, goroutineNum)
	remainingRoutineNum := goroutineNum

	// goroutine退出时计数器减1，关闭out通道，上层停止遍历
	existFunc := func() {
		atomic.AddInt32(&remainingRoutineNum, -1)
		if remainingRoutineNum == 0 {
			close(out)
		}
	}

	for i := int32(0); i < goroutineNum; i++ {
		go safeGoWithCtx(ctx, func() {
			for {
				select {
				// 监听Done通道，若parent ctx发出退出信号，则退出
				case <-ctx.Done():
					logger.Errorf(ctx, "Worker cancel, err: %v", ctx.Err())
					existFunc()
					return
				// 从channel中取req，若没取到，则说明消费完了
				// 否则执行f
				case req, ok := <-in:
					if !ok {
						existFunc()
						return
					}
					reqList, _ := req.([]interface{})
					resp := f(ctx, reqList)
					out <- resp
				}
			}
		})
	}
	return out
}

type WorkerFunc func(ctx context.Context, reqList []interface{}) DoWithGoroutineResp

type DoWithGoroutineResp struct {
	Resp []interface{}
	Err  error
}

func DoWithGoroutine(ctx context.Context, f WorkerFunc, req []interface{}, limit int) ([]interface{}, error) {
	defer logger.TimeCost(ctx)()
	if len(req) == 0 {
		return make([]interface{}, 0), nil
	}

	// 根据入参req数量和切片长度创建channel
	in := make(chan interface{}, len(req)/limit+1)
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	// 最大协程数为10
	goroutineNum := int32(math.Min(float64(len(req)/limit+1), 10))

	// 启动worker
	out := worker(cancelCtx, f, in, goroutineNum)

	// 往channel中写入req
	for i := 0; i < len(req); i += limit {
		end := int(math.Min(float64(i+limit), float64(len(req))))
		in <- req[i:end]
	}
	close(in)

	resp := make([]interface{}, 0, len(req))
	for r := range out {
		rData := r.(DoWithGoroutineResp)
		if rData.Err != nil {
			cancelFunc()
			return nil, errors.Wrap(rData.Err, "[DoWithGoroutine] func run failed")
		}
		resp = append(resp, rData.Resp...)
	}
	cancelFunc()
	return resp, nil
}
