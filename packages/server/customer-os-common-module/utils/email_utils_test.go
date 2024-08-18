package utils

import (
	"testing"
)

func TestGetReadableNameFromEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{"HR case", "hr@hp.com", "Hr"},
		{"Info case", "into@hp.com", "Into"},
		{"Dot-separated case", "alex.smith@hp.com", "Alex Smith"},
		{"Underscore-separated case", "alex_baba@apple.com", "Alex Baba"},
		{"Hyphen-separated case", "m-g@msft.com", "M G"},
		{"Mixed separator case", "john.doe_smith-jones@example.com", "John Doe Smith Jones"},
		{"Single letter case", "a@b.com", "A"},
		{"Empty username case", "@domain.com", ""},
		{"Empty email case", "", ""},
		{"No domain case", "username", "Username"},
		{"Multiple @ symbols", "user@name@domain.com", "User"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetReadableNameFromEmail(tt.email)
			if result != tt.expected {
				t.Errorf("GetReadableNameFromEmail(%q) = %q, want %q", tt.email, result, tt.expected)
			}
		})
	}
}
