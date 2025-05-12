package channels

import (
	"context"

	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/aws"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/rabbitmq"
)

type Channel struct {
	Client      *rabbitmq.Client
	AWSApp      *aws.AWSApplication
	ChannelName string
	QueueName   string
	RoutingKeys []string
	ctx         context.Context
	cancelFunc  context.CancelFunc
}

func NewChannel(client *rabbitmq.Client, app *aws.AWSApplication, channelName string, queueName string, routingKeys []string) *Channel {
	ctx, cancel := context.WithCancel(context.Background())

	return &Channel{
		Client:      client,
		AWSApp:      app,
		ChannelName: channelName,
		QueueName:   queueName,
		RoutingKeys: routingKeys,
		ctx:         ctx,
		cancelFunc:  cancel,
	}
}
