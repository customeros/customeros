package utils

import (
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func ConvertJsonbToStringSlice(input pgtype.JSONB) []string {
	var mappedInput []any
	err := input.AssignTo(&mappedInput)
	if err != nil {
		logrus.Error("Error converting jsonb to string slice: ", err)
		return []string{}
	}
	var output []string
	for _, c := range mappedInput {
		if _, ok := c.(string); ok {
			output = append(output, c.(string))
		} else if _, ok := c.(int64); ok {
			item := strconv.FormatInt(c.(int64), 10)
			output = append(output, item)
		} else if _, ok := c.(float64); ok {
			item := strconv.FormatFloat(c.(float64), 'f', 0, 64)
			output = append(output, item)
		}
	}
	return output
}

func GetUniqueElements[V comparable](input []V) []V {
	if len(input) == 0 {
		return []V{}
	}

	m := make(map[V]bool)
	for _, item := range input {
		if _, ok := m[item]; !ok {
			m[item] = true
		}
	}

	result := make([]V, 0, len(m))
	for key := range m {
		result = append(result, key)
	}
	return result
}

func LowercaseStrings(arr []string) {
	for i, s := range arr {
		arr[i] = strings.ToLower(s)
	}
}

func FloatToString(f *float64) string {
	if f == nil {
		return ""
	}
	return fmt.Sprintf("%f", *f)
}
