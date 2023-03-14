package logger

import (
	"context"
	"runtime"
	"time"
)

func TimeCost(ctx context.Context) func() {
	return TimeCostWithKV(ctx, nil)
}

func TimeCostWithKV(ctx context.Context, kv map[string]interface{}) func() {
	start := time.Now()
	pc, _, _, _ := runtime.Caller(2)
	method := runtime.FuncForPC(pc).Name()
	if kv == nil {
		kv = map[string]interface{}{}
	}
	kv[ContextLogKey_Method] = method
	InfofWithKV(ctx, kv, "[time_cost] %v enter", method)
	return func() {
		cost := time.Since(start).Milliseconds()
		kv[ContextLogKey_TimeCost] = cost
		InfofWithKV(ctx, kv, "[time_cost] %v done, cost: %v", method, cost)
	}
}
