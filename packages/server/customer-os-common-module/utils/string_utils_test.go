package utils

import (
	"testing"
)

func TestExtractFirstPart(t *testing.T) {
	// Positive test case
	str := "Hello, World!"
	delimiter := ","
	expected := "Hello"
	result := ExtractFirstPart(str, delimiter)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Negative test case
	str = "Hello, World!"
	delimiter = ":"
	expected = "Hello, World!"
	result = ExtractFirstPart(str, delimiter)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Edge case: empty string
	str = ""
	delimiter = ","
	expected = ""
	result = ExtractFirstPart(str, delimiter)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestCapitalizeAllParts(t *testing.T) {
	testCases := []struct {
		input      string
		delimiters []string
		expected   string
	}{
		{"america/los_angeles", []string{"/", "_"}, "America/Los_Angeles"},
		{"EUROPE/london", []string{"/"}, "Europe/London"},
		{"hello_world", []string{"_"}, "Hello_World"},
		{"hello_world", []string{}, "Hello_world"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := CapitalizeAllParts(tc.input, tc.delimiters)
			if result != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, result)
			}
		})
	}
}
