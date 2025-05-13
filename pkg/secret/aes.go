package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

// AESEncrypt encrypts plaintext using the provided secret key with AES-GCM mode (which is recommended for authenticated encryption) and returns base64-encoded ciphertext. The secret key should be 16, 24, or 32 bytes long. If there's an error during encryption, it returns empty string and error.
func AESEncrypt(plaintext string, secretKey string) (string, error) {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt decrypts ciphertext using the provided secret key with AES-GCM mode (which is recommended for authenticated encryption) and returns plaintext as string. The secret key should be 16, 24, or 32 bytes long. If there's an error during encryption, it returns empty string and error.
func AESDecrypt(ciphertext string, secretKey string) (string, error) {
	test, _ := base64.URLEncoding.DecodeString(ciphertext)

	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := test[:nonceSize], string(test)[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
