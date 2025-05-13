package awsgw

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthRequest struct {
	ConnectionID string `json:"connection_id"`
	APIKey       string `json:"api_key"`
	Signature    string `json:"signature"`
	Timestamp    string `json:"timestamp"`
}

type AuthResponse struct {
	Event string
	Data  interface{}
}

func (s *AWSGatewayService) AuthWebsocket(c echo.Context) error {
	var req AuthRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	// Verify
	resp, err := s.AuthAPI.Validate("POST", req.APIKey, "")
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to validate API key"})
	}

	// TODO: Save connection ID and user ID
	msg := map[string]string{"message": "authenticated successfully"}
	r := AuthResponse{Event: "auth", Data: msg}
	s.App.PostToConnection(c.Request().Context(), req.ConnectionID, r)

	return c.JSON(http.StatusOK, resp)
}
