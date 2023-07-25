package data

import "testing"

func TestAdjustOrganizationMarket(t *testing.T) {

	testCases := []struct {
		newValue      string
		previousValue string
		expected      string
	}{
		{"", "", ""},
		{"B2B", "", "B2B"},
		{"b2c", "B2B", "B2C"},
		{"Marketplace", "B2B", "Marketplace"},
		{"B2B,B2C,MarketPlace", "B2C", "B2B"},
		{"B2C,MarketPlace", "B2B", "B2C"},
		{"Some Value", "B2B", "B2B"},
		{"Some Value", "Some Other Value", "Some Value"},
		{"Some Value", "", "Some Value"},
		{"Can be b2b", "B2C", "B2B"},
	}

	for _, tc := range testCases {
		result := AdjustOrganizationMarket(tc.newValue, tc.previousValue)
		if result != tc.expected {
			t.Errorf("Expected %s but got %s for input (%s, %s)",
				tc.expected, result, tc.newValue, tc.previousValue)
		}
	}

}
