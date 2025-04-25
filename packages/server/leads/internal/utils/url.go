package utils

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type URLComponents struct {
	URL          string
	Protocol     string
	Host         string
	Path         string
	PathSegments []string
	QueryRaw     string
	QueryParams  map[string][]string
	Anchor       string
}

func ParseURL(rawURL string) (*URLComponents, error) {
	if rawURL == "" {
		return nil, nil
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Extract path parts
	pathParts := []string{}
	if parsedURL.Path != "" {
		// Split by '/' and filter out empty strings
		for _, part := range strings.Split(parsedURL.Path, "/") {
			if part != "" {
				pathParts = append(pathParts, part)
			}
		}
	}

	return &URLComponents{
		URL:          rawURL,
		Protocol:     parsedURL.Scheme,
		Host:         parsedURL.Hostname(),
		Path:         parsedURL.Path,
		PathSegments: pathParts,
		QueryRaw:     parsedURL.RawQuery,
		QueryParams:  parsedURL.Query(),
		Anchor:       parsedURL.Fragment,
	}, nil
}

func SortUrlsByLength(urls []string) []string {
	sort.Slice(urls, func(i, j int) bool {
		// If lengths are different, sort by length
		if len(urls[i]) != len(urls[j]) {
			return len(urls[i]) < len(urls[j])
		}
		// If lengths are equal, sort alphabetically
		return urls[i] < urls[j]
	})
	return urls
}

func StripUrlToBasePath(s string) string {
	if s == "" {
		return ""
	}
	clean := strings.Split(s, "?")[0]
	return NormalizeUrlPath(clean)
}

func NormalizeUrlPath(s string) string {
	if s == "" {
		return ""
	}
	clean := strings.TrimPrefix(s, "https://")
	clean = strings.TrimPrefix(clean, "http://")
	clean = strings.TrimPrefix(clean, "www.")
	return strings.Trim(clean, "/")
}

func IsSuspiciousURL(url string) bool {
	suspiciousPatterns := []string{
		"oastify.com",
		"burpcollaborator.net",
		"interactsh.com",
		"/zws",
		"xss.ht",
		"ngrok.io",
	}

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(url, pattern) {
			return true
		}
	}

	// Check for very long random-looking subdomains
	parts := strings.Split(url, "/")
	if len(parts) >= 3 {
		domain := parts[2]
		if len(domain) > 30 && URLContainsRandomString(domain) {
			return true
		}
	}

	return false
}

func URLContainsRandomString(s string) bool {
	if strings.Contains(s, "google") {
		return false
	}

	// Count letters, numbers, and special characters
	letters := 0
	numbers := 0
	special := 0

	for _, char := range s {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			letters++
		} else if char >= '0' && char <= '9' {
			numbers++
		} else {
			special++
		}
	}

	// If it has a mix of characters and is sufficiently long
	if letters > 0 && numbers > 0 && len(s) > 20 {
		// Check if it has a high entropy (i.e., appears random)
		// Simple heuristic: Over 20 chars with mixed numbers and letters
		return true
	}

	return false
}
