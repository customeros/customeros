package utils

import (
	"net/url"
	"strings"
	"time"
)

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

// Deprecated use from common module
func ExtractDomainFromUrl(inputURL string) string {
	// Prepend "http://" if the URL doesn't start with a scheme
	if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		inputURL = "http://" + inputURL
	}

	// Parse the URL
	u, err := url.Parse(inputURL)
	if err != nil {
		return ""
	}

	// Extract and return the hostname (domain)
	domain := u.Hostname()

	// Remove "www." if it exists
	if strings.HasPrefix(domain, "www.") {
		domain = domain[4:] // Remove the first 4 characters ("www.")
	}

	return strings.ToLower(domain)
}
