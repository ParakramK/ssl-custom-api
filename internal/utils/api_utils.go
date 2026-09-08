package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateKey() (string, error) {
	b := make([]byte, 32) // 256 bits

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate API key: %w", err)
	}

	return "ssl_" + base64.RawURLEncoding.EncodeToString(b), nil
}
