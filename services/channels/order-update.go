package channels

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/rabbitmq"
)

func (ch *Channel) StartOrderUpdateChannel() error {
	return ch.Client.StartConsumer(ch.ctx, ch.QueueName, ch.RoutingKeys, ch.OrderUpdateMessageHandler)
}

func (ch *Channel) OrderUpdateMessageHandler(msg rabbitmq.Message) error {
	res := NewSuccessMessage(ch.ChannelName, msg)
	// Handle incoming messages from RabbitMQ
	_, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return fmt.Errorf("cannot marshaling message")
	}

	return nil
}

func (ch *Channel) Stop() {
	if ch.cancelFunc != nil {
		ch.cancelFunc()
	}
}
