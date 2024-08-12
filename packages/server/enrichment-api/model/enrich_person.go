package model

import postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

type EnrichPersonRequest struct {
	Email       string `json:"email"`
	LinkedinUrl string `json:"linkedinUrl"`
}

type EnrichPersonResponse struct {
	Status   string              `json:"status"`
	Message  string              `json:"message,omitempty"`
	RecordId uint64              `json:"recordId,omitempty"`
	Data     *EnrichedPersonData `json:"data,omitempty"`
}

type EnrichedPersonData struct {
	PersonProfile *postgresEntity.ScrapInPersonResponse `json:"scrapinPersonProfile,omitempty"`
}
