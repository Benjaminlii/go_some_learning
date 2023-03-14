package rocketmq

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"

	"github.com/Benjaminlii/go_some_learning/common/config"
	"github.com/Benjaminlii/go_some_learning/common/logger"
)

const (
	RetryTime = 3
	SleepTime = time.Millisecond * 10
)

var (
	mqConfig     *config.RocketMQConfig
	Producer     rocketmq.Producer
	PushConsumer rocketmq.PushConsumer
)

func InitRocketMQ() {
	mqConfig = config.GetConfig().RocketMQ
	if mqConfig == nil {
		panic("[InitRocketMQ] failed, mq config is nil")
	}
	var err error
	if Producer, err = rocketmq.NewProducer(producer.WithNameServer(mqConfig.NameServers)); err != nil {
		panic(fmt.Sprintf("[InitRocketMQ] init rocket mq producer err:%v", err))
	}
	if err = Producer.Start(); err != nil {
		panic(fmt.Sprintf("[InitRocketMQ] producer mq start err:%v", err))
	}
	if PushConsumer, err = rocketmq.NewPushConsumer(consumer.WithNameServer(mqConfig.NameServers)); err != nil {
		panic(fmt.Sprintf("[InitRocketMQ] init rocket mq push consumer err:%v", err))
	}
}

func SendMessage(ctx context.Context, body string, tags []string) (err error) {
	msg := primitive.NewMessage(mqConfig.Topic, []byte(body))
	if len(tags) > 0 {
		msg = msg.WithTag(strings.Join(tags, "||"))
	}
	for i := 0; i < RetryTime; i++ {
		res, err := Producer.SendSync(ctx, msg)
		if err != nil {
			time.Sleep(SleepTime)
			continue
		}
		logger.Infof(ctx, "[RocketMQ][SendMessage] send success, msgID:%s", res.MsgID)
		break
	}
	return
}

// 发送延时消息
func SendDelayMessage(ctx context.Context, body string, tags []string, delayLevel int) (err error) {
	msg := primitive.NewMessage(mqConfig.Topic, []byte(body))
	msg.WithDelayTimeLevel(delayLevel)
	if len(tags) > 0 {
		msg = msg.WithTag(strings.Join(tags, "||"))
	}
	for i := 0; i < RetryTime; i++ {
		res, err := Producer.SendSync(ctx, msg)
		if err != nil {
			time.Sleep(SleepTime)
			continue
		}
		logger.Infof(ctx, "[RocketMQ][SendDelayMessage] send success, msgID:%s", res.MsgID)
		break
	}
	return

}
