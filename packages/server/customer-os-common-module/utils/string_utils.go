package utils

import (
	"bytes"
	"crypto/rand"
	"github.com/forPelevin/gomoji"
	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math/big"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	charsetAlphaNumeric      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLowerAlphaNumeric = "abcdefghijklmnopqrstuvwxyz0123456789"
	charsetLowerAlpha        = "abcdefghijklmnopqrstuvwxyz"
	charsetSpecial           = "!@#$%^&*()-_=+[]{}<>?"
)

func GenerateLowerAlpha(length int) string {
	if length < 1 {
		return ""
	}
	bytes := make([]byte, length)
	for i := range bytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charsetLowerAlpha))))
		if err != nil {
			panic(err)
		}
		bytes[i] = charsetLowerAlpha[num.Int64()]
	}
	return string(bytes)
}

func GenerateKey(length int, includeSpecial bool) string {
	alphaNumericLength := length
	if includeSpecial {
		alphaNumericLength--
	}
	if alphaNumericLength < 1 {
		return ""
	}
	bytes := make([]byte, alphaNumericLength)
	for i := range bytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charsetLowerAlphaNumeric))))
		if err != nil {
			panic(err)
		}
		bytes[i] = charsetLowerAlphaNumeric[num.Int64()]
	}
	if includeSpecial {
		specialCharIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charsetSpecial))))
		if err != nil {
			panic(err)
		}
		bytes = append(bytes, charsetSpecial[specialCharIndex.Int64()])
	}
	return string(bytes)
}

func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	for i := range bytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charsetAlphaNumeric))))
		if err != nil {
			panic(err)
		}
		bytes[i] = charsetAlphaNumeric[num.Int64()]
	}
	return string(bytes)
}

func JoinNonEmpty(delimiter string, strs ...string) string {
	var nonEmptyStrs []string
	for _, s := range strs {
		if len(s) > 0 {
			nonEmptyStrs = append(nonEmptyStrs, s)
		}
	}
	return strings.Join(nonEmptyStrs, delimiter)
}

func StringOrEmpty(key *string) string {
	if key != nil {
		return *key
	}
	return ""
}

func StringFirstNonEmpty(strs ...string) string {
	for _, s := range strs {
		if len(s) > 0 {
			return s
		}
	}
	return ""
}

func StringPtrFirstNonEmpty(strs ...*string) string {
	for _, s := range strs {
		if s == nil {
			continue
		}
		if *s != "" {
			return *s
		}
		if len(*s) > 0 {
			return *s
		}
	}
	return ""
}

func NewUUIDIfEmpty(str string) string {
	if strings.TrimSpace(str) == "" {
		return uuid.New().String()
	}
	return strings.TrimSpace(str)
}

func ExtractFirstPart(str, delimiter string) string {
	// Find the first delimiter
	delimiterIndex := strings.Index(str, delimiter)
	if delimiterIndex == -1 {
		// No delimiter found, return the whole string
		return str
	}
	// Extract the first part
	firstPart := str[:delimiterIndex]
	return firstPart
}

func CapitalizeAllParts(str string, delimiters []string) string {
	if len(delimiters) == 0 {
		titleCase := cases.Title(language.Und)
		return titleCase.String(str)
	}

	for _, delimiter := range delimiters {
		str = capitalizeParts(str, delimiter)
	}
	return str
}

func capitalizeParts(input, delimiter string) string {
	parts := strings.Split(input, delimiter)

	// Create a title casing for capitalizing the words
	titleCase := cases.Title(language.Und)

	// Capitalize the first letter of each word
	for i, part := range parts {
		parts[i] = titleCase.String(part)
	}

	// Join the parts back together
	capitalized := strings.Join(parts, delimiter)

	return capitalized
}

func SliceToString(slice []string) string {
	return strings.Join(slice, ",")
}

func StringToSlice(str string) []string {
	if str == "" {
		return []string{}
	}
	return strings.Split(str, ",")
}

func NormalizeString(input string) string {
	replacements := map[rune]string{
		'é': "e",
		'è': "e",
		'ê': "e",
		'ë': "e",
		'à': "a",
		'â': "a",
		'ô': "o",
		'ö': "o",
		'û': "u",
		'ü': "u",
		'ï': "i",
		'î': "i",
		'ç': "c",
		'ñ': "n",
	}

	// remove emojis
	input = gomoji.RemoveEmojis(input)

	// replace special characters
	var result strings.Builder
	for _, r := range input {
		if replacement, ok := replacements[unicode.ToLower(r)]; ok {
			if unicode.IsUpper(r) {
				result.WriteString(strings.ToUpper(replacement))
			} else {
				result.WriteString(replacement)
			}
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

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

func ToCamelCase(input string) string {
	if input == "" {
		return ""
	}
	return strings.ToUpper(string(input[0])) + strings.ToLower(input[1:])
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
