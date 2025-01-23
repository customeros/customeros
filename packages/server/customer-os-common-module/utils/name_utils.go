package utils

import (
	"github.com/forPelevin/gomoji"
	"strings"
)

// SplitFullName splits a fullName into (firstName, lastName).
// If there's only one token, it goes to firstName and lastName is empty.
// If there are multiple tokens, the first token is the firstName and
// the rest of the tokens form the lastName, joined by a space.
//
// Rejoining firstName + " " + lastName should yield the original name (minus extra spacing).
func SplitFullName(fullName string) (string, string) {
	// Trim leading/trailing whitespace
	s := strings.TrimSpace(fullName)
	if s == "" {
		return "", ""
	}

	// Split on whitespace
	tokens := strings.Fields(s)

	// If only one token, treat it as firstName, leave lastName empty
	if len(tokens) == 1 {
		return tokens[0], ""
	}

	// Otherwise, first token is firstName; remaining tokens become lastName
	firstName := tokens[0]
	lastName := strings.Join(tokens[1:], " ")

	return firstName, lastName
}

func CleanName(input string) string {
	if input == "" {
		return input
	}
	output := input
	output = strings.ReplaceAll(output, "\n", "")
	output = strings.ReplaceAll(output, "\r", "")
	output = strings.ReplaceAll(output, "\t", "")
	output = strings.ReplaceAll(output, "\v", "")
	output = strings.ReplaceAll(output, "\f", "")

	// replace special characters
	specialChars := []string{
		"®", "™", "©", "℠", "&amp;", "&", "@", "!", "?", "*", "#",
		"(", ")", "[", "]", "{", "}", "|", "\\", "/", "+", "=",
	}
	for _, char := range specialChars {
		output = strings.ReplaceAll(output, char, "")
	}

	// Remove emojis
	output = gomoji.RemoveEmojis(output)

	// Remove non-ASCII characters
	output = SanitizeUTF8(output)

	// Capitalize the first letter of each word
	output = CapitalizeAllParts(output, []string{" "})

	return output
}
