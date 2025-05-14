package channels

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Gaoey/poc-aws-websocket-gateway.git/services/awsgw"
)

func (ch *Channel) StartOrderUpdateChannel() error {
	return ch.Client.StartConsumer(ch.ctx, ch.QueueName, ch.RoutingKeys, ch.OrderUpdateMessageHandler)
}

func (ch *Channel) OrderUpdateMessageHandler(msg interface{}) error {
	var msgData map[string]interface{}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error marshaling:", err)
		return nil
	}

	if err := json.Unmarshal(msgBytes, &msgData); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return nil
	}

	var userId string
	var ok bool
	if userId, ok = msgData["user_id"].(string); !ok {
		log.Printf("user_id not found in message")
		return nil
	}

	key := awsgw.GetWSChannelKey(ch.ChannelName, userId)
	members, err := ch.Redis.SMembers(ch.ctx, key)
	if err != nil {
		log.Printf("Error getting members from Redis: %v", err)
		return nil
	}

	res := NewSuccessMessage(ch.ChannelName, msg)
	for _, member := range members {
		err := ch.AWSApp.PostToConnection(ch.ctx, member, res)
		if err != nil {
			log.Printf("Error sending message to connection %s: %v", member, err)
			continue
		}
	}

	return nil
}

func (ch *Channel) Stop() {
	if ch.cancelFunc != nil {
		ch.cancelFunc()
	}
}
