package entity

import "testing"

func TestGetNamesFromString(t *testing.T) {
	c := ContactEntity{}

	tests := []struct {
		input         string
		expectedFirst string
		expectedLast  string
	}{
		{"john", "John", ""},
		{"john doe", "John", "Doe"},
		{"john.doe😃", "John", "Doe"},
		{"john-doe", "John", "Doe"},
		{"john_doe", "John", "Doe"},
		{"john+doe", "John", "Doe"},
		{"☺️😄😂😇john=DOE", "John", "Doe"},
		{"john,doe", "John", "Doe"},
		{"Alice Wonderland", "Alice", "Wonderland"},
		{"michael-jordan", "Michael", "Jordan"},
		{"steve.jobs", "Steve", "Jobs"},
		{"jane-doe smith", "Jane", "Doe smith"},
		{"  john doe  ", "John", "Doe"}, // Extra spaces
		{"JOHN", "John", ""},            // Single name input
	}

	for _, test := range tests {
		firstName, lastName := c.GetNamesFromString(test.input)
		if firstName != test.expectedFirst {
			t.Errorf("For input '%s', expected first name '%s' but got '%s'", test.input, test.expectedFirst, firstName)
		}
		if lastName != test.expectedLast {
			t.Errorf("For input '%s', expected last name '%s' but got '%s'", test.input, test.expectedLast, lastName)
		}
	}
}
