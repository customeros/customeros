package utils

import (
	"bytes"
	"unicode/utf8"
)

func SanitizeUTF8(input string) string {
	// Check if the string is already valid UTF-8
	if utf8.ValidString(input) {
		return input
	}

	// If not, replace invalid UTF-8 sequences
	validString := bytes.Buffer{}
	for i := 0; i < len(input); {
		r, size := utf8.DecodeRuneInString(input[i:])
		if r == utf8.RuneError && size == 1 {
			// Replace invalid byte with a placeholder, e.g., '?'
			validString.WriteRune('?')
			i++
		} else {
			validString.WriteRune(r)
			i += size
		}
	}
	return validString.String()
}
