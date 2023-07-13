package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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

func Hmac(body []byte, key []byte) *string {
	// Create a new HMAC hasher with the desired hash function and secret key
	hasher := hmac.New(sha256.New, key)

	// Write the message to the hasher
	hasher.Write(body)

	// Calculate the HMAC value
	hmacValue := hasher.Sum(nil)

	// Convert the HMAC value to a hexadecimal string representation
	hmacString := hex.EncodeToString(hmacValue)
	return &hmacString
}
