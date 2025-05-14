package awsgw

import (
	"net/http"

	"github.com/Gaoey/poc-aws-websocket-gateway.git/domain"
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

const WS_AUTH_KEY = "ws-auth"

func GetWSAuthKey(userID string) string {
	return WS_AUTH_KEY + ":" + userID
}

func (s *AWSGatewayService) AuthWebsocket(c echo.Context) error {
	var req AuthRequest
	if err := c.Bind(&req); err != nil {
		resp := domain.WSResponse{
			Event: "auth",
			Data:  map[string]string{"error": "Invalid request", "message": err.Error()},
		}
		return c.JSON(400, resp)
	}

	// Verify
	resp, err := s.AuthAPI.Validate("POST", req.Payload.Data.APIKey, req.Payload.Data.Timestamp, req.Payload.Data.Signature, "")
	if err != nil {
		resp := domain.WSResponse{
			Event: "auth",
			Data:  map[string]string{"error": "Failed to validate API key", "message": err.Error()},
		}
		return c.JSON(500, resp)
	}

	// Save connection ID and user ID to redis
	key := GetWSAuthKey(req.ConnectionID)
	data := map[string]interface{}{
		"user_id": resp.UserID,
		"is_auth": true,
	}
	if err := s.Redis.SetHashData(c.Request().Context(), key, data); err != nil {
		resp := domain.WSResponse{
			Event: "auth",
			Data:  map[string]string{"error": "Failed to save connection ID", "message": err.Error()},
		}
		return c.JSON(500, resp)
	}

	msg := map[string]string{"message": "authenticated successfully"}
	r := domain.WSResponse{Event: "auth", Data: msg}

	return c.JSON(http.StatusOK, r)
}
