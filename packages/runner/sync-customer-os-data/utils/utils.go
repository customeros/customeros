package utils

import "time"

const customLayout1 = "2006-01-02 15:04:05"
const customLayout2 = "2022-11-07T22:03:16.000+0000"

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
	if err != nil {
		return nil, err
	}

	t, err = time.Parse(customLayout2, input)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
