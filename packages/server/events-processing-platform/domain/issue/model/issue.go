package model

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type Issue struct {
	ID                     string                  `json:"id"`
	Tenant                 string                  `json:"tenant"`
	Subject                string                  `json:"subject"`
	Description            string                  `json:"description"`
	Status                 string                  `json:"status"`
	Priority               string                  `json:"priority"`
	ReportedByOrganization string                  `json:"loggedOrganizationIds,omitempty"`
	Source                 cmnmod.Source           `json:"source"`
	ExternalSystems        []cmnmod.ExternalSystem `json:"externalSystem"`
	CreatedAt              time.Time               `json:"createdAt,omitempty"`
	UpdatedAt              time.Time               `json:"updatedAt,omitempty"`
}
