package awsgw

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *AWSGatewayService) ConntectWebSocket(c echo.Context) error {
	msg := map[string]string{"message": "connected websocket successfully"}
	return c.JSON(http.StatusOK, msg)
}

func (s *AWSGatewayService) DisconnectWebSocket(c echo.Context) error {
	// TODO: Clear here
	return c.JSON(http.StatusOK, map[string]string{
		"message": "WebSocket connection closed",
	})
}
