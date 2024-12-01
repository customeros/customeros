package service

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type Integration string

// Add all supported integrations here, and also in validIntegrations below

const (
	IntegrationCalCom     Integration = "calcom"
	IntegrationFathom     Integration = "fathom"
	IntegrationGrain      Integration = "grain"
	IntegrationHubspot    Integration = "hubspot"
	IntegrationPostmark   Integration = "postmark"
	IntegrationSalesforce Integration = "salesforce"
)

var ValidIntegrations = func() map[string]Integration {
	integrations := []Integration{
		IntegrationCalCom,
		IntegrationFathom,
		IntegrationGrain,
		IntegrationHubspot,
		IntegrationPostmark,
		IntegrationSalesforce,
	}

	m := make(map[string]Integration)
	for _, i := range integrations {
		m[string(i)] = i
	}
	return m
}()

func (i Integration) String() string {
	return string(i)
}

func (i Integration) IntegrationID(rotationCount int64) string {
	// Add rotation count to string being hashed
	input := fmt.Sprintf("%s:%d", i.String(), rotationCount)
	return utils.GenerateHashId(input, 12)
}
