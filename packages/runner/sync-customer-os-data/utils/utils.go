package utils

import "time"

const customLayout = "2006-01-02 15:04:05"

func UnmarshalDateTime(input string) (*time.Time, error) {
	if input == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, input)
	if err == nil {
		// Parsed as RFC3339
		return &t, nil
	}

	// Try custom layout
	t, err = time.Parse(customLayout, input)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
