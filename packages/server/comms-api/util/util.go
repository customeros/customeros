package util

import (
	"fmt"
	"net/mail"
	"strings"
)

func EnsureRfcId(id string) string {
	if !strings.HasPrefix(id, "<") {
		id = fmt.Sprintf("<%s>", id)
	}
	return id
}

func EnsureRfcIds(to []string) []string {
	var result []string
	for _, id := range to {
		result = append(result, EnsureRfcId(id))
	}
	return result
}

func toStringArr(from []*mail.Address) []string {
	var to []string
	for _, a := range from {
		to = append(to, a.Address)
	}
	return to
}
func FirstNotEmpty(input ...string) *string {
	for _, item := range input {
		if item != "" {
			return &item
		}
	}
	return nil
}
