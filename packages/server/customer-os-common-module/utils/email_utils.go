package utils

import (
	"strings"
	"unicode"
)

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

func ExtractDomainFromEmail(email string) string {
	emailParts := strings.Split(email, "@")
	if len(emailParts) != 2 {
		return ""
	}
	return ExtractDomain(emailParts[1])
}
