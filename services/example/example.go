package example

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/rabbitmq"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/pkg/secret"
	"github.com/labstack/echo/v4"
)

type BodyPayload struct {
	RoutingKey string      `json:"routing_key"`
	Message    interface{} `json:"message"`
}

type ExampleHandler struct {
	Client *rabbitmq.Client
}

func NewExampleHandler(client *rabbitmq.Client) *ExampleHandler {
	return &ExampleHandler{
		Client: client,
	}
}

func (h *ExampleHandler) PublishMessage(c echo.Context) error {
	// Parse request body
	var payload BodyPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Publish message to RabbitMQ
	err := h.Client.Publish(c.Request().Context(), payload.RoutingKey, payload.Message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to publish message",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "message published successfully",
	})
}

type SignaturePayload struct {
	APIKey string `json:"api_key"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   string `json:"body"`
}

func (h *ExampleHandler) GetSignature(c echo.Context) error {
	// Parse request body
	var payload SignaturePayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})

	}

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	key := os.Getenv("SECRET_KEY")
	p := timestamp + payload.Method + payload.Path + payload.Body
	signature := secret.SignHmacSha256([]byte(p), key)

	data := map[string]string{
		"api_key":   payload.APIKey,
		"signature": signature,
		"timestamp": timestamp,
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"event": "auth",
		"data":  data,
	})
}
