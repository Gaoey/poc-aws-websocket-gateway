package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Gaoey/poc-aws-websocket-gateway.git/pkg/secret"
)

type AuthAPI struct {
	URL string
}

func NewAuthAPI(endpoint string) *AuthAPI {
	return &AuthAPI{}
}

type ValidationResponse struct {
	UserID     string           `json:"user_id"`
	Permission APIKeyPermission `json:"permission"`
} //@name ValidationResponse

type APIKeyPermission struct {
	IsTrade    bool `json:"is_trade"`
	IsWithdraw bool `json:"is_withdraw"`
	IsDeposit  bool `json:"is_deposit"`
	IsInternal bool `json:"is_internal"`
	IsAirdrop  bool `json:"is_airdrop"`
} //@name APIKeyPermission

func (a *AuthAPI) Validate(method, apiKey string, body string) (*ValidationResponse, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	key := os.Getenv("SECRET_KEY")
	endpoint := "api/v1/validation"
	p := timestamp + method + endpoint + body
	signature := secret.SignHmacSha256([]byte(p), key)

	req, err := http.NewRequest(method, a.URL+endpoint, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-BTK-APIKEY", apiKey)
	req.Header.Set("X-BTK-SIGN", signature)
	req.Header.Set("X-BTK-TIMESTAMP", timestamp)
	req.Header.Set("X-BTK-METHOD", method)
	req.Header.Set("X-BTK-PATH", endpoint)

	fmt.Printf("headers: %#v\n", req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to validate API key: %s", resp.Status)
	}
	defer resp.Body.Close()
	var validationResponse ValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&validationResponse); err != nil {
		return nil, fmt.Errorf("failed to decode validation response: %v", err)
	}

	return &validationResponse, nil
}
