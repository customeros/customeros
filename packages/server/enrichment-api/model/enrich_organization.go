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

type EnrichOrganizationResponse struct {
	Status              string                         `json:"status"`
	Message             string                         `json:"message,omitempty"`
	Success             bool                           `json:"success"`
	PrimaryEnrichSource string                         `json:"primaryEnrichSource"`
	Data                EnrichOrganizationResponseData `json:"data"`
}

type EnrichOrganizationResponseData struct {
	Name             string                                 `json:"name"`
	Domain           string                                 `json:"domain"`
	ShortDescription string                                 `json:"description"`
	LongDescription  string                                 `json:"longDescription"`
	Website          string                                 `json:"website"`
	Employees        int64                                  `json:"employees"`
	FoundedYear      int64                                  `json:"foundedYear"`
	Public           *bool                                  `json:"public,omitempty"`
	Logos            []string                               `json:"logos"`
	Icons            []string                               `json:"icons"`
	Industry         string                                 `json:"industry"`
	Socials          []EnrichOrganizationResponseSocial     `json:"socials"`
	Location         EnrichOrganizationResponseDataLocation `json:"location"`
}

type EnrichOrganizationResponseSocial struct {
	Url   string `json:"url"`
	Alias string `json:"alias"`
	Id    string `json:"id"`
}

type EnrichOrganizationResponseDataLocation struct {
	IsHeadquarter *bool  `json:"isHeadquarter"`
	Country       string `json:"country"`
	CountryCodeA2 string `json:"countryCodeA2"`
	CountryCodeA3 string `json:"countryCodeA3"`
	Locality      string `json:"locality"`
	Region        string `json:"region"`
	PostalCode    string `json:"postalCode"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2"`
}

func (l EnrichOrganizationResponseDataLocation) IsEmpty() bool {
	return l.Country == "" && l.CountryCodeA2 == "" && l.Locality == "" && l.Region == ""
}

type EnrichOrganizationScrapinResponse struct {
	Status            string                              `json:"status"`
	Message           string                              `json:"message,omitempty"`
	RecordId          uint64                              `json:"recordId,omitempty"`
	OrganizationFound bool                                `json:"organizationFound"`
	Data              *postgresEntity.ScrapInResponseBody `json:"data,omitempty"`
}
