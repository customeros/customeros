package data

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
)

// Market type constants
const (
	B2B         = "B2B"
	B2C         = "B2C"
	Marketplace = "Marketplace"
)

func AdjustOrganizationMarket(newValue, previousValue string) string {
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
	} else if utils.Contains([]string{B2B, B2C, Marketplace}, previousValue) {
		return previousValue
	}
	return newValue
}
