package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"strings"
)

func GetSignature(input, key string) string {
	keyForSign := []byte(key)
	h := hmac.New(sha1.New, keyForSign)
	h.Write([]byte(input))
	hash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	log.Printf("hash len=%d", len(hash))

	for ; len(hash)%4 != 0; hash = hash + "=" {
		log.Printf("Padding hash len=%d", len(hash))
	}
	hash = strings.Replace(hash, " ", "+", -1)
	return hash
}
