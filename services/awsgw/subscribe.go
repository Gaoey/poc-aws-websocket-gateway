package awsgw

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type Channel string

const (
	OrderUpdateChannel Channel = "order-update"
)

type SubscribePayload struct {
	Event   string  `json:"event"`
	Channel Channel `json:"channel"`
}

type SubscribeRequest struct {
	ConnectionID string           `json:"connection_id"`
	Payload      SubscribePayload `json:"payload"`
}

type SubscribeResponse struct {
	event string
	data  interface{}
}

func (s *AWSGatewayService) SubscribeChannel(c echo.Context) error {
	var req SubscribeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	// TODO: Save connectionID to redis
	fmt.Printf("\n\nReceived subscription request: %#v\n\n", req)
	msg := map[string]string{"message": "Subscribed to channel successfully"}
	s.App.PostToConnection(c.Request().Context(), req.ConnectionID, msg)

	return c.JSON(200, SubscribeResponse{event: req.Payload.Event, data: msg})
}
