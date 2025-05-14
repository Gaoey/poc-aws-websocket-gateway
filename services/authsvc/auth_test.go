package authsvc_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/Gaoey/poc-aws-websocket-gateway.git/pkg/secret"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/services/awsgw"
	"github.com/joho/godotenv"
)

func TestValidate(t *testing.T) {
	// Load .env from project root
	projectRoot := os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		// Fallback to relative path from test file
		_, filename, _, _ := runtime.Caller(0)
		projectRoot = filepath.Join(filepath.Dir(filename), "../..")
	}
	err := godotenv.Load(filepath.Join(projectRoot, ".env"))
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	method := "POST"
	path := "/validation"
	body := ""
	apiKey := "3c501338-0121-431b-9415-291c8b50f1ec"
	userSecretKey := "DRjDa4gfY7ZG1aW_ZZ_Vf9JN9hFEIxVUs3ze9sb_SCkIMSjrXaSucMm9gLKiE1GGl8t4FiwlvTODTJ5vR2s5PK2ee6pRVvKPUNNmO_YXCW7vdjuEMIGxxSJL_37E81LMwdw243Fcqx9X9g59lguyIg6Nr0jbsqDt"
	secretKey := os.Getenv("SECRET_KEY")

	decrypt, err := secret.AESDecrypt(userSecretKey, secretKey)
	if err != nil {
		t.Fatalf("Failed to decrypt secret key: %v", err)
	}

	// Create timestamp and signature
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	payload := timestamp + method + path + body
	signature := secret.SignHmacSha256([]byte(payload), decrypt)

	data := awsgw.AuthData{
		APIKey:    apiKey,
		Signature: signature,
		Timestamp: timestamp,
	}
	msg := map[string]interface{}{
		"event": "auth",
		"data":  data,
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message to JSON: %v", err)
	}

	t.Logf("\n\n%s\n\n", string(jsonMsg))
	fmt.Printf("Test data - API Key: %s, Timestamp: %s\n", data.APIKey, data.Timestamp)

	// If you want to actually validate something in this test:
	if data.APIKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, data.APIKey)
	}
}
