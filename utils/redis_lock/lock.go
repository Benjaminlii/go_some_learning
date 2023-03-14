package redislock

import (
	"context"
	"errors"
	"time"

	"github.com/Benjaminlii/async_task/biz/driver/redis"
	"github.com/Benjaminlii/go_some_learning/common/logger"
)

type Lock struct {
	key     string
	timeout time.Duration
}

// GetSpinLock 基于redis的自旋锁
// locktimeout：锁的超时时间
// spinTimeout
func GetSpinLock(ctx context.Context, key string, lockTimeout, spinTimeout time.Duration) (*Lock, error) {
	l := &Lock{key: key, timeout: lockTimeout}
	lockSuccuss := false
	ch := make(chan struct{})
	spinTimeoutCh := time.NewTimer(spinTimeout).C
	go func() {
		for {
			select {
			case <-spinTimeoutCh:
				close(ch)
				return
			default:
				err := l.Try()
				if err == nil {
					lockSuccuss = true
					close(ch)
					return
				}
				logger.Warnf(ctx, "[GetSpinLock] lock.Try failed. key: %s, err: %s",
					key, err.Error())
				time.Sleep(20 * time.Millisecond)
			}
		}
	}()
	<-ch
	if !lockSuccuss {
		return nil, errors.New("[GetSpinLock] try get lock timeout")
	}
	return l, nil
}

func (l *Lock) Try() error {
	ok, err := redis.Redis().SetNX(l.key, 1, l.timeout).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("try lock failed")
	}
	return nil
}

func (l *Lock) Release(ctx context.Context) error {
	_, err := redis.Redis().Del(l.key).Result()
	logger.Infof(ctx, "[Release] unlock Key %s", l.key)
	return err
}
