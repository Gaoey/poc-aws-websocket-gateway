package authsvc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/redis"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/pkg/secret"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthService struct {
	Redis *redis.RedisHandler
}

func NewAuthService(redis *redis.RedisHandler) *AuthService {
	return &AuthService{
		Redis: redis,
	}
}

type CreateAPIKeyRequest struct {
	Label string `json:"label"`
}

type APIKeyData struct {
	UserID      string           `json:"user_id"`
	LabelName   string           `json:"label_name"`
	APIKey      string           `json:"api_key"`
	SecretKey   string           `json:"api_secret"`
	Permission  APIKeyPermission `json:"permission"`
	IPWhitelist []string         `json:"ip_whitelist"`
	CreatedAt   int64            `json:"created_at"`
	UpdatedAt   int64            `json:"updated_at"`
}

type APIKeyPermission struct {
	IsTrade    bool `json:"is_trade"`
	IsWithdraw bool `json:"is_withdraw"`
	IsDeposit  bool `json:"is_deposit"`
	IsInternal bool `json:"is_internal"`
	IsAirdrop  bool `json:"is_airdrop"`
}

func (s *AuthService) CreateAPIKey(c echo.Context) error {
	secretKey := os.Getenv("SECRET_KEY")
	userID := c.Param("user_id")
	if userID == "" {
		return echo.NewHTTPError(400, "user_id is required")
	}

	var err error
	var req CreateAPIKeyRequest
	err = json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}

	newApiKey := uuid.New().String()
	newSecret, err := secret.GenSecret(100)
	if err != nil {
		return echo.NewHTTPError(500, "Failed to generate secret key")
	}
	secretCiphertext, err := secret.AESEncrypt(newSecret, secretKey)
	if err != nil {
		return echo.NewHTTPError(500, "Failed to encrypt secret key")
	}

	currentTime := time.Now().UnixMilli()
	data := APIKeyData{
		UserID:    userID,
		LabelName: req.Label,
		APIKey:    newApiKey,
		SecretKey: secretCiphertext,
		Permission: APIKeyPermission{
			IsTrade:    false,
			IsWithdraw: false,
			IsDeposit:  false,
			IsInternal: false,
			IsAirdrop:  false,
		},
		IPWhitelist: []string{},
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return echo.NewHTTPError(500, "Failed to encrypt secret key")
	}

	// save to redis
	prefix := os.Getenv("REDIS_KEY_PREFIX")
	key := prefix + ":" + newApiKey
	err = s.Redis.SetJSONData(context.Background(), key, ".", jsonData)
	if err != nil {
		return echo.NewHTTPError(500, "Failed to save API key to Redis")
	}

	return c.JSON(200, map[string]interface{}{
		"data": data,
	})
}

type ValidationData struct {
	APIKey    string `json:"X-BTK-APIKEY"`
	TimeStamp int64  `json:"X-BTK-TIMESTAMP"`
	Signature string `json:"X-BTK-SIGN"`
	Method    string `json:"X-BTK-METHOD"`
	Path      string `json:"X_BTK_PATH"`
	IP        string `json:"X_BTK_IP"`
	Body      string `json:"body"`
}

type ValidationResponse struct {
	UserID     string           `json:"user_id"`
	Permission APIKeyPermission `json:"permission"`
} //@name ValidationResponse

type HeaderRequest struct {
	APIKey    string `header:"X-BTK-APIKEY" validate:"required"`
	Sign      string `header:"X-BTK-SIGN" validate:"required"`
	Timestamp string `header:"X-BTK-TIMESTAMP" validate:"required,numeric,len=13"`
	Method    string `header:"X-BTK-METHOD" validate:"required"`
	Path      string `header:"X-BTK-PATH" validate:"required"`
	IP        string `header:"X-BTK-IP" validate:"required,ipv4"`
}

func (s *AuthService) ValidateAPIKey(c echo.Context) error {
	headers := new(HeaderRequest)
	if err := (&echo.DefaultBinder{}).BindHeaders(c, headers); err != nil {
		return echo.NewHTTPError(400, "Invalid request headers")
	}

	// get data from redis
	var data APIKeyData
	prefix := os.Getenv("REDIS_KEY_PREFIX")
	index := fmt.Sprintf("%s:%s", prefix, headers.APIKey)
	d, err := s.Redis.GetJSONData(context.Background(), index, ".")
	if err != nil {
		return echo.NewHTTPError(500, "Failed to get API key from Redis")
	}

	err = json.Unmarshal([]byte(d.(string)), &data)
	if err != nil {
		return echo.NewHTTPError(500, "Failed to unmarshal API key data")
	}

	// Get Secret key
	secretKey := os.Getenv("SECRET_KEY")

	decrypt, err := secret.AESDecrypt(data.SecretKey, secretKey)
	if err != nil {
		return echo.NewHTTPError(500, "Failed to decrypt secret key")
	}

	bb, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(500, "Failed to read request body")
	}

	payload := headers.Timestamp + headers.Method + headers.Path + string(bb)
	expectedSignature := secret.SignHmacSha256([]byte(payload), decrypt)

	if expectedSignature != headers.Sign {
		return echo.NewHTTPError(401, "Invalid signature")
	}

	return c.JSON(200, map[string]interface{}{
		"user_id":     data.UserID,
		"permissions": data.Permission})
}
