package model

import (
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"strings"
)

type EnrichOrganizationRequest struct {
	Domain      string `json:"domain"`
	LinkedinUrl string `json:"linkedinUrl"`
}

func (e *EnrichOrganizationRequest) Normalize() {
	e.LinkedinUrl = strings.TrimSpace(e.LinkedinUrl)
	e.Domain = strings.TrimSpace(e.Domain)
}

type EnrichOrganizationScrapinResponse struct {
	Status            string                              `json:"status"`
	Message           string                              `json:"message,omitempty"`
	RecordId          uint64                              `json:"recordId,omitempty"`
	OrganizationFound bool                                `json:"organizationFound"`
	Data              *postgresEntity.ScrapInResponseBody `json:"data,omitempty"`
}
