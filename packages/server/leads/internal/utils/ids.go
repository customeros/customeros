package utils

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateNanoIDWithPrefix(prefix string, length int) string {
	alphabet := "abcdefghijklmnopqrstuvwxyz0123456789"
	id, err := gonanoid.Generate(alphabet, length)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s_%s", prefix, id)
}

func GenerateNanoID(length int) string {
	alphabet := "abcdefghijklmnopqrstuvwxyz0123456789"
	id, err := gonanoid.Generate(alphabet, length)
	if err != nil {
		panic(err)
	}
	return id
}
