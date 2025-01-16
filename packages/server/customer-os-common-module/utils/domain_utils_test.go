package utils

import "testing"

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "HTTP scheme + www",
			input:    "http://www.example.com",
			expected: "example.com",
		},
		{
			name:     "HTTPS scheme + www",
			input:    "https://www.example.com",
			expected: "example.com",
		},
		{
			name:     "Invalid URL - no dots",
			input:    "invalidurl",
			expected: "",
		},
		{
			name:     "Domain without scheme",
			input:    "example.com",
			expected: "example.com",
		},
		{
			name:     "Subdomains - no scheme",
			input:    "hu.example.com",
			expected: "example.com",
		},
		{
			name:     "Subdomains - co.uk",
			input:    "ro.example.co.uk",
			expected: "example.co.uk",
		},
		{
			name:     "Deep subdomains - co.uk",
			input:    "home.iasi.ro.example.co.uk",
			expected: "example.co.uk",
		},
		{
			name:     "Subdomains - cop.ro",
			input:    "example.cop.ro",
			expected: "cop.ro",
		},
		{
			name:     "Subdomains - co.ro",
			input:    "example.co.ro",
			expected: "example.co.ro",
		},
		{
			name:     "Domain - co.uk",
			input:    "example.co.uk",
			expected: "example.co.uk",
		},
		{
			name:     "Wrong TLD",
			input:    "http://www.example.stupidme",
			expected: "",
		},
		{
			name:     "URL with path",
			input:    "http://www.example.com/main/final?param1=true&param2=1",
			expected: "example.com",
		},
		{
			name:     "Email-like input",
			input:    "alex@opeline.com.py",
			expected: "opeline.com.py",
		},
		{
			name:     "Empty String",
			input:    "",
			expected: "",
		},
		{
			name:     "Single TLD only",
			input:    "com",
			expected: "",
		},
		{
			name:     "Trailing slash",
			input:    "example.com/",
			expected: "example.com",
		},
		{
			name:     "Uppercase Domain",
			input:    "EXAMPLE.COM",
			expected: "example.com",
		},
		{
			name:     "Mixed case + scheme + path",
			input:    "HTTP://WWW.EXAMPLE.COM/TestPage",
			expected: "example.com",
		},
		{
			name:     "Leading + Trailing whitespace",
			input:    "   www.example.com  ",
			expected: "example.com",
		},
		{
			name:     "Domain + port",
			input:    "http://example.com:8080",
			expected: "example.com",
		},
		// Note: This next test only makes sense if your code supports IDNA or punycode. If not, it's good future-proofing.
		{
			name:     "Internationalized Domain Name (punycode)",
			input:    "http://xn--exmple-cua.com",
			expected: "xn--exmple-cua.com",
			// or the Unicode-decoded version if you do IDNA decoding
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ExtractDomain(tc.input)
			if actual != tc.expected {
				t.Errorf("Input: %q | Expected: %q, got: %q", tc.input, tc.expected, actual)
			}
		})
	}
}
