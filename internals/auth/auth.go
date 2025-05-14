package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthAPI struct {
	URL string
}

func NewAuthAPI(endpoint string) *AuthAPI {
	return &AuthAPI{
		URL: endpoint,
	}
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

func (a *AuthAPI) Validate(method, apiKey string, ts string, signature string, body string) (*ValidationResponse, error) {
	endpoint := "/validation"

	req, err := http.NewRequest(method, a.URL+endpoint, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-BTK-APIKEY", apiKey)
	req.Header.Set("X-BTK-SIGN", signature)
	req.Header.Set("X-BTK-TIMESTAMP", ts)
	req.Header.Set("X-BTK-METHOD", method)
	req.Header.Set("X-BTK-PATH", endpoint)

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
