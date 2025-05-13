package secret

import (
	"strings"

	"github.com/google/uuid"
)

func GenSecret(length int) (string, error) {
	// Generate UUID round 1
	uuid1, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	// Generate UUID round 2
	uuid2, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	// Generate random string
	randomStr, err := GenerateRandomString(length - 72)
	if err != nil {
		return "", err
	}

	// Remove dash ("-") from UUID.
	r1 := strings.ReplaceAll(uuid1.String(), "-", "")
	r2 := strings.ReplaceAll(uuid2.String(), "-", "")

	// Concatenate uuid and random string.
	secretKey := r1 + r2 + randomStr
	return secretKey, nil
}
