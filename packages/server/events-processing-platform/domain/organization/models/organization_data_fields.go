package models

import (
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"time"
)

type OrganizationDataFields struct {
	Name              string
	Description       string
	Website           string
	Industry          string
	SubIndustry       string
	IndustryGroup     string
	TargetAudience    string
	ValueProposition  string
	IsPublic          bool
	Employees         int64
	Market            string
	LastFundingRound  string
	LastFundingAmount string
	SlackChannelLink  string
}

type OrganizationFields struct {
	ID                     string
	Tenant                 string
	IgnoreEmptyFields      bool
	OrganizationDataFields OrganizationDataFields
	Source                 commonModels.Source
	CreatedAt              *time.Time
	UpdatedAt              *time.Time
}
