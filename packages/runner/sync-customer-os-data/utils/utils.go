package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

const customLayout1 = "2006-01-02 15:04:05"
const customLayout2 = "2006-01-02T15:04:05.000-0700"

func UnmarshalDateTime(input string) (*time.Time, error) {
	if input == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, input)
	if err == nil {
		// Parsed as RFC3339
		return &t, nil
	}

	// Try custom layouts
	t, err = time.Parse(customLayout1, input)
	if err == nil {
		return &t, nil
	}

	t, err = time.Parse(customLayout2, input)
	if err == nil {
		return &t, nil
	}

	return nil, errors.New(fmt.Sprintf("cannot parse input as date time %s", input))
}
