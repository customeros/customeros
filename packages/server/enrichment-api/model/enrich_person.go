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
	Domain      string `json:"domain"`
}

func (e *EnrichPersonRequest) Normalize() {
	e.Email = strings.TrimSpace(e.Email)
	e.LinkedinUrl = strings.TrimSpace(e.LinkedinUrl)
	e.FirstName = strings.TrimSpace(e.FirstName)
	e.LastName = strings.TrimSpace(e.LastName)
	e.Domain = strings.TrimSpace(e.Domain)
}

type EnrichPersonScrapinResponse struct {
	Status      string              `json:"status"`
	Message     string              `json:"message,omitempty"`
	RecordId    uint64              `json:"recordId,omitempty"`
	PersonFound bool                `json:"personFound"`
	Data        *EnrichedPersonData `json:"data,omitempty"` // TODO replace directly with *postgresEntity.ScrapInResponseBody
}

type EnrichedPersonData struct {
	PersonProfile *postgresEntity.ScrapInResponseBody `json:"scrapinPersonProfile,omitempty"`
}
