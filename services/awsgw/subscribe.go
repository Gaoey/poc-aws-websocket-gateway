package awsgw

import (
	"fmt"

	"github.com/Gaoey/poc-aws-websocket-gateway.git/domain"
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

func GetWSChannelKey(channel string, userID string) string {
	return fmt.Sprintf("ws-channel:%s:%s", channel, userID)
}

func (s *AWSGatewayService) SubscribeChannel(c echo.Context) error {
	var req SubscribeRequest
	if err := c.Bind(&req); err != nil {
		resp := domain.WSResponse{
			Event: "subscribe",
			Data:  map[string]string{"error": "Invalid request", "message": err.Error()},
		}
		return c.JSON(400, resp)
	}

	// CHECK is authenticated
	data, err := s.Redis.GetHashData(c.Request().Context(), GetWSAuthKey(req.ConnectionID))
	if err != nil {
		resp := domain.WSResponse{
			Event: "subscribe",
			Data:  map[string]string{"error": "Failed to get user data", "message": err.Error()},
		}
		return c.JSON(500, resp)
	}

	if data["is_auth"] != "1" {
		resp := domain.WSResponse{
			Event: "subscribe",
			Data:  map[string]string{"error": "Unauthorized", "message": "User is not authenticated"},
		}
		return c.JSON(401, resp)
	}

	UserID := data["user_id"]
	key := GetWSChannelKey(string(req.Payload.Channel), UserID)
	isMem, err := s.Redis.SIsMember(c.Request().Context(), key, req.ConnectionID)
	if err != nil {
		resp := domain.WSResponse{
			Event: "subscribe",
			Data:  map[string]string{"error": "Failed to check channel membership", "message": err.Error()},
		}
		return c.JSON(500, resp)
	}

	if !isMem {
		err = s.Redis.SAdd(c.Request().Context(), key, req.ConnectionID)
		if err != nil {
			resp := domain.WSResponse{
				Event: "subscribe",
				Data:  map[string]string{"error": "Failed to subscribe to channel", "message": err.Error()},
			}
			return c.JSON(500, resp)
		}
	}

	msg := map[string]string{"message": "Subscribed to channel successfully"}
	resp := domain.WSResponse{Event: req.Payload.Event, Data: msg}

	return c.JSON(200, resp)
}
