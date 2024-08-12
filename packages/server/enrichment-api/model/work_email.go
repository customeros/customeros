package model

import postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

type FindWorkEmailRequest struct {
	LinkedinUrl   string `json:"linkedinUrl"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	CompanyDomain string `json:"companyDomain"`
}

type FindWorkEmailResponse struct {
	Status   string                                    `json:"status"`
	Message  string                                    `json:"message,omitempty"`
	RecordId string                                    `json:"recordId,omitempty"`
	Data     *postgresentity.BetterContactResponseBody `json:"data,omitempty"`
}
