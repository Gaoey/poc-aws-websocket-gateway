package awsgw

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ConnectPayload struct {
	ConnectionID string `json:"connection_id"`
}

type ConnectRequest struct {
	ConnectionID string           `json:"connection_id"`
	Payload      SubscribePayload `json:"payload"`
}

type ConnectResponse struct {
	Event string
	Data  interface{}
}

func (s *AWSGatewayService) ConntectWebSocket(c echo.Context) error {
	var req ConnectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	// TODO: Verify here
	msg := map[string]string{"message": "connected websocket successfully"}
	resp := ConnectResponse{Event: "auth", Data: msg}
	s.App.PostToConnection(c.Request().Context(), req.ConnectionID, resp)

	return c.JSON(http.StatusOK, resp)
}

func (s *AWSGatewayService) DisconnectWebSocket(c echo.Context) error {
	// TODO: Clear here
	return c.JSON(http.StatusOK, map[string]string{
		"message": "WebSocket connection closed",
	})
}
