package utils

import (
	"github.com/google/uuid"
	"strings"
)

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
	return str
}
