package utils

import (
	"github.com/google/uuid"
	"strings"
)

func JoinNonEmpty(delimiter string, strs ...string) string {
	var nonEmptyStrs []string
	for _, s := range strs {
		if len(s) > 0 {
			nonEmptyStrs = append(nonEmptyStrs, s)
		}
	}
	return strings.Join(nonEmptyStrs, delimiter)
}

func StringFirstNonEmpty(strs ...string) string {
	for _, s := range strs {
		if len(s) > 0 {
			return s
		}
	}
	return ""
}

func NewUUIDIfEmpty(str string) string {
	if strings.TrimSpace(str) == "" {
		return uuid.New().String()
	}
	return strings.TrimSpace(str)
}

func ExtractFirstPart(str, delimiter string) string {
	// Find the first delimiter
	delimiterIndex := strings.Index(str, delimiter)
	if delimiterIndex == -1 {
		// No delimiter found, return the whole string
		return str
	}
	// Extract the first part
	firstPart := str[:delimiterIndex]
	return firstPart
}
