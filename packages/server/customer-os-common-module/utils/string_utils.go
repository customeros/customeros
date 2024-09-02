package utils

import (
	"crypto/rand"
	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math/big"
	"strings"
	"unicode"
)

const (
	charset                  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLowerAlphaNumeric = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func GenerateKey(length int) string {
	bytes := make([]byte, length)
	for i := range bytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charsetLowerAlphaNumeric))))
		if err != nil {
			panic(err)
		}
		bytes[i] = charsetLowerAlphaNumeric[num.Int64()]
	}
	return string(bytes)
}

func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	for i := range bytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		bytes[i] = charset[num.Int64()]
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

func NormalizeString(s string) string {
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

	var result strings.Builder
	for _, r := range s {
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
