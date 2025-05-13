package secret

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// This function computes the HMAC-SHA256 signature for the given payload using the provided secret key and returns the signature as a hexadecimal string.
func SignHmacSha256(payload []byte, secret string) string {
	// create hash function (hmac sha256) using the provided key
	hm := hmac.New(sha256.New, []byte(secret))
	// write payload
	hm.Write(payload)
	//convert to hex
	return hex.EncodeToString(hm.Sum(nil))
}
