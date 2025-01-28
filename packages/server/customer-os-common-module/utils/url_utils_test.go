package utils

import "testing"

func TestCleanUrlBasePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Plain domain",
			input:    "example.com",
			expected: "example.com",
		},
		{
			name:     "Strip leading slash",
			input:    "/example.com",
			expected: "example.com",
		},
		{
			name:     "Strip https://",
			input:    "https://example.com",
			expected: "example.com",
		},
		{
			name:     "Strip http://",
			input:    "http://example.com",
			expected: "example.com",
		},
		{
			name:     "Strip www.",
			input:    "www.example.com",
			expected: "example.com",
		},
		{
			name:     "Remove query string",
			input:    "https://www.example.com/path?query=param",
			expected: "example.com/path",
		},
		{
			name:     "Trailing slash",
			input:    "https://example.com/",
			expected: "example.com",
		},
		{
			name:     "Path with trailing slash",
			input:    "http://www.example.com/foo/bar/",
			expected: "example.com/foo/bar",
		},
		{
			name:     "Complex path + query + trailing slash",
			input:    "https://www.example.org/foo/bar/baz/?q=test",
			expected: "example.org/foo/bar/baz",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CleanUrlBasePath(tc.input)
			if got != tc.expected {
				t.Errorf("CleanUrlBasePath(%q) = %q; want %q", tc.input, got, tc.expected)
			}
		})
	}
}
