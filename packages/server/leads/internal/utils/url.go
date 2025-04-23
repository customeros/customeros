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
