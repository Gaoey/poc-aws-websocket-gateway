package secret

import (
	"crypto/rand"
	"math/big"
)

var (
	CHARSET   string = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"
	DEFAULT_N int    = 30
)

func GenerateRandomString(n int) (string, error) {
	//  Uses the default number if n < 0, in case the SECRET_LENGTH environment variable was not provided.
	if n < 0 {
		n = DEFAULT_N
	}

	// Generate random string
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(CHARSET))))
		if err != nil {
			return "", err
		}
		ret[i] = CHARSET[num.Int64()]
	}

	return string(ret), nil
}
