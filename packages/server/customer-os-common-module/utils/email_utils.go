package utils

import (
	"fmt"
	"strings"
	"unicode"
)

func EnsureEmailRfcId(id string) string {
	if !strings.HasPrefix(id, "<") {
		id = fmt.Sprintf("<%s>", id)
	}
	return id
}

func EnsureEmailRfcIds(to []string) []string {
	if to == nil {
		return nil
	}
	var result []string
	for _, id := range to {
		result = append(result, EnsureEmailRfcId(id))
	}
	return result
}

func GetReadableNameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 0 {
		return ""
	}

	username := parts[0]

	// Split the username by ., -, and _
	words := strings.FieldsFunc(username, func(r rune) bool {
		return r == '.' || r == '-' || r == '_'
	})

	// Capitalize first letter of each word and join with spaces
	for i, word := range words {
		if len(word) > 0 {
			words[i] = string(unicode.ToUpper(rune(word[0]))) + word[1:]
		}
	}

	return strings.Join(words, " ")
}
