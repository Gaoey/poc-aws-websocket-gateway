package awsgw

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *AWSGatewayService) ConntectWebSocket(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "WebSocket connection established",
	})
}

func (s *AWSGatewayService) DisconnectWebSocket(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "WebSocket connection closed",
	})
}
