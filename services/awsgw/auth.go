package awsgw

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthData struct {
	APIKey    string `json:"api_key"`
	Signature string `json:"signature"`
	Timestamp string `json:"timestamp"`
}
type AuthPayload struct {
	Event string   `json:"event"`
	Data  AuthData `json:"data"`
}

type AuthRequest struct {
	ConnectionID string      `json:"connection_id"`
	Payload      AuthPayload `json:"payload"`
}

type AuthResponse struct {
	Event string
	Data  interface{}
}

func (s *AWSGatewayService) AuthWebsocket(c echo.Context) error {
	var req AuthRequest
	if err := c.Bind(&req); err != nil {
		resp := AuthResponse{
			Event: "auth",
			Data:  map[string]string{"error": "Invalid request", "message": err.Error()},
		}
		return c.JSON(400, resp)
	}

	// Verify
	resp, err := s.AuthAPI.Validate("POST", req.Payload.Data.APIKey, req.Payload.Data.Timestamp, req.Payload.Data.Signature, "")
	if err != nil {
		resp := AuthResponse{
			Event: "auth",
			Data:  map[string]string{"error": "Failed to validate API key", "message": err.Error()},
		}
		return c.JSON(500, resp)
	}

	// TODO: Save connection ID and user ID
	msg := map[string]string{"message": "authenticated successfully"}
	r := AuthResponse{Event: "auth", Data: msg}
	s.App.PostToConnection(c.Request().Context(), req.ConnectionID, r)

	return c.JSON(http.StatusOK, resp)
}
