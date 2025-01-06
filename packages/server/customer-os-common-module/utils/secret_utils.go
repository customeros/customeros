package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateSecret generates a cryptographically secure random string of the given byte size,
// encoded in URL-safe Base64 without padding.
func GenerateSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
