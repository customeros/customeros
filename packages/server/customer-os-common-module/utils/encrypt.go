package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func Encrypt(plaintext, encodedKey string) (string, string, error) {
	// Decode the base64 encoded key
	secretKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return "", "", errors.New("invalid encryption key")
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return hex.EncodeToString(ciphertext[aes.BlockSize:]), hex.EncodeToString(iv), nil
}

func Decrypt(encryptedHex, ivHex string, encodedKey string) (string, error) {
	// Decode the base64 encoded key
	secretKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return "", errors.New("invalid encryption key")
	}

	encrypted, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}

	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)

	return string(encrypted), nil
}
