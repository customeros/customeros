package utils

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
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

func GenerateUUID() (string, error) {
	uuidObj, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return uuidObj.String(), nil
}
