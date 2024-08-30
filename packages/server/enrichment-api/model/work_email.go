package model

import postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

type FindWorkEmailRequest struct {
	LinkedinUrl       string `json:"linkedinUrl"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	CompanyName       string `json:"companyName"`
	CompanyDomain     string `json:"companyDomain"`
	EnrichPhoneNumber bool   `json:"enrichPhoneNumber"`
}

type FindWorkEmailResponse struct {
	Status                 string                                    `json:"status"`
	Message                string                                    `json:"message,omitempty"`
	RecordId               string                                    `json:"recordId,omitempty"`
	BetterContactRequestId string                                    `json:"betterContactRequestId,omitempty"`
	Data                   *postgresentity.BetterContactResponseBody `json:"data,omitempty"`
}
