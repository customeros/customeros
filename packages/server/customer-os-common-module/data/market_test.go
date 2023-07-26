package data

import "testing"

func TestAdjustOrganizationMarket(t *testing.T) {

	testCases := []struct {
		newValue string
		expected string
	}{
		{"", ""},
		{"B2B", "B2B"},
		{"b2c", "B2C"},
		{"Marketplace", "Marketplace"},
		{"B2B,B2C,MarketPlace", "B2B"},
		{"B2C,MarketPlace", "B2C"},
		{"Some Value", "Some Value"},
		{"Can be b2b", "B2B"},
	}

	for _, tc := range testCases {
		result := AdjustOrganizationMarket(tc.newValue)
		if result != tc.expected {
			t.Errorf("Expected %s but got %s for input (%s)",
				tc.expected, result, tc.newValue)
		}
	}

}
