package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHashId(s string, length int) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))[:length]
}
