package models

import (
	"fmt"
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"time"
)

type Social struct {
	PlatformName string `json:"platformName"`
	Url          string `json:"url"`
}

type Organization struct {
	ID                string                             `json:"id"`
	Name              string                             `json:"name"`
	Description       string                             `json:"description"`
	Website           string                             `json:"website"`
	Industry          string                             `json:"industry"`
	SubIndustry       string                             `json:"subIndustry"`
	IndustryGroup     string                             `json:"industryGroup"`
	TargetAudience    string                             `json:"targetAudience"`
	ValueProposition  string                             `json:"valueProposition"`
	IsPublic          bool                               `json:"isPublic"`
	Employees         int64                              `json:"employees"`
	Market            string                             `json:"market"`
	LastFundingRound  string                             `json:"lastFundingRound"`
	LastFundingAmount string                             `json:"lastFundingAmount"`
	Source            common_models.Source               `json:"source"`
	CreatedAt         time.Time                          `json:"createdAt"`
	UpdatedAt         time.Time                          `json:"updatedAt"`
	PhoneNumbers      map[string]OrganizationPhoneNumber `json:"phoneNumbers"`
	Emails            map[string]OrganizationEmail       `json:"emails"`
	Domains           []string                           `json:"domains,omitempty"`
	Socials           map[string]Social                  `json:"socials,omitempty"`
}

type OrganizationPhoneNumber struct {
	Primary bool   `json:"primary"`
	Label   string `json:"label"`
}

type OrganizationEmail struct {
	Primary bool   `json:"primary"`
	Label   string `json:"label"`
}

func (o *Organization) String() string {
	return fmt.Sprintf("Organization{ID: %s, Name: %s, Description: %s, Website: %s, Industry: %s, IsPublic: %t, Source: %s, CreatedAt: %s, UpdatedAt: %s}", o.ID, o.Name, o.Description, o.Website, o.Industry, o.IsPublic, o.Source, o.CreatedAt, o.UpdatedAt)
}
