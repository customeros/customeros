package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

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
