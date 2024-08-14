package model

import (
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"strings"
)

type EnrichPersonRequest struct {
	Email       string `json:"email"`
	LinkedinUrl string `json:"linkedinUrl"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
}

func (e *EnrichPersonRequest) Normalize() {
	e.Email = strings.TrimSpace(e.Email)
	e.LinkedinUrl = strings.TrimSpace(e.LinkedinUrl)
	e.FirstName = strings.TrimSpace(e.FirstName)
	e.LastName = strings.TrimSpace(e.LastName)
}

type EnrichPersonResponse struct {
	Status      string              `json:"status"`
	Message     string              `json:"message,omitempty"`
	RecordId    uint64              `json:"recordId,omitempty"`
	PersonFound bool                `json:"personFound"`
	Data        *EnrichedPersonData `json:"data,omitempty"`
}

type EnrichedPersonData struct {
	PersonProfile *postgresEntity.ScrapInPersonResponse `json:"scrapinPersonProfile,omitempty"`
}
