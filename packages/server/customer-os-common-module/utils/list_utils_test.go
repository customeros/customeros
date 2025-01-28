package utils

import (
	"reflect"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Empty slice",
			input:    []string{},
			expected: nil,
		},
		{
			name:     "No duplicates",
			input:    []string{"apple", "banana", "orange"},
			expected: []string{"apple", "banana", "orange"},
		},
		{
			name:     "All duplicates",
			input:    []string{"apple", "apple", "apple"},
			expected: []string{"apple"},
		},
		{
			name:     "Some duplicates",
			input:    []string{"apple", "banana", "apple", "orange", "banana"},
			expected: []string{"apple", "banana", "orange"},
		},
		{
			name:     "Mixed case duplicates",
			input:    []string{"Apple", "apple", "Banana", "BANANA", "banana"},
			expected: []string{"Apple", "apple", "Banana", "BANANA", "banana"},
		},
		{
			name:     "Duplicates not contiguous",
			input:    []string{"a", "b", "c", "a", "d", "b", "e"},
			expected: []string{"a", "b", "c", "d", "e"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := RemoveDuplicates(tc.input)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("RemoveDuplicates(%v) = %v; want %v", tc.input, got, tc.expected)
			}
		})
	}
}
