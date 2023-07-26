package utils

import (
	"fmt"
	"strings"
)

func EnsureEmailRfcId(id string) string {
	if !strings.HasPrefix(id, "<") {
		id = fmt.Sprintf("<%s>", id)
	}
	return id
}

func EnsureEmailRfcIds(to []string) []string {
	var result []string
	for _, id := range to {
		result = append(result, EnsureEmailRfcId(id))
	}
	return result
}
