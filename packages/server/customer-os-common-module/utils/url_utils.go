package utils

import "strings"

func CleanUrlBasePath(s string) string {
	if s == "" {
		return ""
	}
	clean := strings.Split(s, "?")[0]
	clean = strings.TrimPrefix(clean, "https://")
	clean = strings.TrimPrefix(clean, "http://")
	clean = strings.TrimPrefix(clean, "www.")
	return strings.Trim(clean, "/")
}
