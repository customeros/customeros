package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/caches"
)

func adjustIndustryValue(inputIndustry string, caches caches.Cache) string {
	if inputIndustry == "" {
		return ""
	}
	if industryValue, ok := caches.GetIndustry(inputIndustry); ok {
		return industryValue
	}
	return ""
}
