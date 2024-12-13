package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateHashId(s string, length int) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))[:length]
}

func GenerateNanoIdWithPrefix(s string, length int) string {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id, err := gonanoid.Generate(alphabet, length)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s-%s", s, id)
}
