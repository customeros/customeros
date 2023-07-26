package data

import (
	"strings"
)

// Market type constants
const (
	B2B         = "B2B"
	B2C         = "B2C"
	Marketplace = "Marketplace"
)

func AdjustOrganizationMarket(newValue string) string {
	if newValue == "" {
		return ""
	}
	marketUpper := strings.ToUpper(newValue)
	if strings.Contains(marketUpper, B2B) {
		return B2B
	} else if strings.Contains(marketUpper, B2C) {
		return B2C
	} else if strings.Contains(marketUpper, "MARKETPLACE") {
		return Marketplace
	}
	return newValue
}
