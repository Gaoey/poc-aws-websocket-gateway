package channels

import (
	"context"

	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/aws"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/rabbitmq"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/redis"
)

type Channel struct {
	Client      *rabbitmq.Client
	AWSApp      *aws.AWSApplication
	Redis       *redis.RedisHandler
	ChannelName string
	QueueName   string
	RoutingKeys []string
	ctx         context.Context
	cancelFunc  context.CancelFunc
}

func NewChannel(client *rabbitmq.Client, app *aws.AWSApplication, redis *redis.RedisHandler, channelName string, queueName string, routingKeys []string) *Channel {
	ctx, cancel := context.WithCancel(context.Background())

	return &Channel{
		Client:      client,
		AWSApp:      app,
		Redis:       redis,
		ChannelName: channelName,
		QueueName:   queueName,
		RoutingKeys: routingKeys,
		ctx:         ctx,
		cancelFunc:  cancel,
	}
}
