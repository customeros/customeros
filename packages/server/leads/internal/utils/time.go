package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	customLayout1 = "2006-01-02 15:04:05"
	customLayout2 = "2006-01-02T15:04:05.000-0700"
	customLayout3 = "2006-01-02T15:04:05-07:00"
	customLayout4 = "Mon, 2 Jan 2006 15:04:05 -0700 (MST)"
	customLayout5 = "Mon, 2 Jan 2006 15:04:05 MST"
	customLayout6 = "Mon, 2 Jan 2006 15:04:05 -0700"
	customLayout7 = "Mon, 2 Jan 2006 15:04:05 +0000 (GMT)"
	customLayout8 = "Mon, 2 Jan 2006 15:04:05 -0700 (MST)"
	customLayout9 = "2 Jan 2006 15:04:05 -0700"
)

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
	customLayouts := []string{customLayout1, customLayout2, customLayout4, customLayout5, customLayout6, customLayout7, customLayout8, customLayout9}

	for _, layout := range customLayouts {
		t, err = time.Parse(layout, input)
		if err == nil {
			return &t, nil
		}
	}
	inputForLayout3 := input
	if !strings.Contains(input, "[UTC]") {
		index := strings.Index(input, "[")
		// If found, strip off the timezone information
		if index != -1 {
			inputForLayout3 = input[:index]
		}
	}
	t, err = time.Parse(customLayout3, inputForLayout3)
	if err == nil {
		return &t, nil
	}

	return nil, errors.New(fmt.Sprintf("cannot parse input as date time %s", input))
}
