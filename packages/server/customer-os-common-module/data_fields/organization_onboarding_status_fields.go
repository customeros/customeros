package data_fields

import neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

type OrganizationOnboardingStatusFields struct {
	Status             *neo4jenum.OnboardingStatus `json:"status,omitempty"`
	Comments           *string                     `json:"comments,omitempty"`
	CausedByContractId *string                     `json:"causedByContractId,omitempty"`
}
