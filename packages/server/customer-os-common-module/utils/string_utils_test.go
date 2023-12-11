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
