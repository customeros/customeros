package utils

import (
	"testing"
)

func TestSplitFullName(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantFirstName string
		wantLastName  string
		wantRejoined  string // expected rejoin of firstName and lastName
	}{
		{
			name:          "Empty string",
			input:         "",
			wantFirstName: "",
			wantLastName:  "",
			wantRejoined:  "",
		},
		{
			name:          "Single token",
			input:         "Cher",
			wantFirstName: "Cher",
			wantLastName:  "",
			wantRejoined:  "Cher",
		},
		{
			name:          "Two tokens",
			input:         "John Smith",
			wantFirstName: "John",
			wantLastName:  "Smith",
			wantRejoined:  "John Smith",
		},
		{
			name:          "Extra spaces",
			input:         "  John   Smith   ",
			wantFirstName: "John",
			wantLastName:  "Smith",
			wantRejoined:  "John Smith",
		},
		{
			name:          "Multiple tokens (3)",
			input:         "John Adam Smith",
			wantFirstName: "John",
			wantLastName:  "Adam Smith",
			wantRejoined:  "John Adam Smith",
		},
		{
			name:          "Multiple tokens (4)",
			input:         "Mr. John Adam Smith",
			wantFirstName: "Mr.",
			wantLastName:  "John Adam Smith",
			wantRejoined:  "Mr. John Adam Smith",
		},
		{
			name:          "With punctuation spacing",
			input:         " Dr.  John   Doe ",
			wantFirstName: "Dr.",
			wantLastName:  "John Doe",
			wantRejoined:  "Dr. John Doe",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotFirst, gotLast := SplitFullName(tc.input)

			// Check firstName/lastName against expectations
			if gotFirst != tc.wantFirstName || gotLast != tc.wantLastName {
				t.Errorf("SplitFullName(%q) = (%q, %q), want (%q, %q)",
					tc.input, gotFirst, gotLast, tc.wantFirstName, tc.wantLastName)
			}

			// Rejoin with no extra space if one part is empty
			var rejoined string
			if gotFirst != "" && gotLast != "" {
				rejoined = gotFirst + " " + gotLast
			} else {
				// If either is empty, just concatenate directly
				rejoined = gotFirst + gotLast
			}

			// Compare with the expected rejoined result
			if rejoined != tc.wantRejoined {
				t.Errorf("Rejoined name = %q, want %q", rejoined, tc.wantRejoined)
			}
		})
	}
}
