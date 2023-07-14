package utils

import (
	"github.com/jackc/pgtype"
	"github.com/sirupsen/logrus"
	"strconv"
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
