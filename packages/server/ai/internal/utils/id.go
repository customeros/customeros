package utils

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateNanoIdWithPrefix(s string, length int) string {
	alphabet := "abcdefghijklmnopqrstuvwxyz0123456789"
	id, err := gonanoid.Generate(alphabet, length)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s_%s", s, id)
}

func GenerateLLMRequestID() string {
	return GenerateNanoIdWithPrefix("llm", 21)
}

func GenerateAPIRequestID() string {
	return GenerateNanoIdWithPrefix("api", 21)
}
