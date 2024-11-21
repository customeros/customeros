package data_fields

import (
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
)

type OrganizationFields struct {
	AppSource          *string                             `json:"appSource,omitempty"`
	Source             *string                             `json:"source,omitempty"`
	ExternalSystem     *model.ExternalSystem               `json:"externalSystem,omitempty"`
	Name               *string                             `json:"name,omitempty"`
	Website            *string                             `json:"website,omitempty"`
	Domains            []string                            `json:"domains,omitempty"`
	Stage              *neo4jenum.OrganizationStage        `json:"stage,omitempty"`
	Relationship       *neo4jenum.OrganizationRelationship `json:"relationship,omitempty"`
	CustomerOsId       *string                             `json:"customerOsId,omitempty"`
	Hide               *bool                               `json:"hide,omitempty"`
	Description        *string                             `json:"description,omitempty"`
	Industry           *string                             `json:"industry,omitempty"`
	SubIndustry        *string                             `json:"subIndustry,omitempty"`
	IndustryGroup      *string                             `json:"industryGroup,omitempty"`
	TargetAudience     *string                             `json:"targetAudience,omitempty"`
	ValueProposition   *string                             `json:"valueProposition,omitempty"`
	LastFundingRound   *string                             `json:"lastFundingRound,omitempty"`
	LastFundingAmount  *string                             `json:"lastFundingAmount,omitempty"`
	ReferenceId        *string                             `json:"referenceId,omitempty"`
	Note               *string                             `json:"note,omitempty"`
	IsPublic           *bool                               `json:"isPublic,omitempty"`
	Employees          *int64                              `json:"employees,omitempty"`
	Market             *string                             `json:"market,omitempty"`
	YearFounded        *int64                              `json:"yearFounded,omitempty"`
	Headquarters       *string                             `json:"headquarters,omitempty"`
	LogoUrl            *string                             `json:"logoUrl,omitempty"`
	IconUrl            *string                             `json:"iconUrl,omitempty"`
	EmployeeGrowthRate *string                             `json:"employeeGrowthRate,omitempty"`
	SlackChannelId     *string                             `json:"slackChannelId,omitempty"`
	LeadSource         *string                             `json:"leadSource,omitempty"`
	IcpFit             *bool                               `json:"icpFit,omitempty"`
	EnrichDomain       *string                             `json:"enrichDomain,omitempty"`
	EnrichSource       *string                             `json:"enrichSource,omitempty"`
	OwnerId            *string                             `json:"ownerId,omitempty"`
	LinkedInUrl        *string                             `json:"linkedInUrl,omitempty"`
}

func (fields OrganizationFields) ExternalSystemAvailable() bool {
	return fields.ExternalSystem != nil && fields.ExternalSystem.Available()
}

func (fields OrganizationFields) GetStageStr() string {
	if fields.Stage != nil {
		return fields.Stage.String()
	}
	return ""
}

func (fields OrganizationFields) GetRelationshipStr() string {
	if fields.Relationship != nil {
		return fields.Relationship.String()
	}
	return ""
}
