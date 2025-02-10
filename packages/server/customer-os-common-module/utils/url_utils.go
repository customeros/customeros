package utils

import (
	"regexp"
	"strings"
)

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

func MatchUrlPattern(pattern, url string) bool {
	// Convert URL pattern wildcards to filepath wildcards
	// * becomes [^/]* (any chars except /)
	// ** becomes .* (any chars including /)

	// First escape any special regex chars
	pattern = regexp.QuoteMeta(pattern)

	// Replace ** with .* (any characters)
	pattern = strings.ReplaceAll(pattern, "\\*\\*", ".*")

	// Replace * with [^/]* (any characters except /)
	pattern = strings.ReplaceAll(pattern, "\\*", "[^/]*")

	// Add anchors and compile regex
	pattern = "^" + pattern + "$"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	return regex.MatchString(url)
}
