package dto

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
)

type UpdateOrganization struct {
	ExternalSystem     *neo4jmodel.ExternalSystem `json:"externalSystem,omitempty"`
	OwnerId            *string                    `json:"ownerId,omitempty"`
	Hide               *bool                      `json:"hide,omitempty"`
	Name               *string                    `json:"name,omitempty"`
	Description        *string                    `json:"description,omitempty"`
	Website            *string                    `json:"website,omitempty"`
	Industry           *string                    `json:"industry,omitempty"`
	SubIndustry        *string                    `json:"subIndustry,omitempty"`
	IndustryGroup      *string                    `json:"industryGroup,omitempty"`
	TargetAudience     *string                    `json:"targetAudience,omitempty"`
	ValueProposition   *string                    `json:"valueProposition,omitempty"`
	IsPublic           *bool                      `json:"isPublic,omitempty"`
	Employees          *int64                     `json:"employees,omitempty"`
	Market             *string                    `json:"market,omitempty"`
	LastFundingRound   *string                    `json:"lastFundingRound,omitempty"`
	LastFundingAmount  *string                    `json:"lastFundingAmount,omitempty"`
	CustomerOsId       *string                    `json:"customerOsId,omitempty"`
	ReferenceId        *string                    `json:"referenceId,omitempty"`
	Note               *string                    `json:"note,omitempty"`
	LogoUrl            *string                    `json:"logoUrl,omitempty"`
	IconUrl            *string                    `json:"iconUrl,omitempty"`
	Headquarters       *string                    `json:"headquarters,omitempty"`
	YearFounded        *int64                     `json:"yearFounded,omitempty"`
	EmployeeGrowthRate *string                    `json:"employeeGrowthRate,omitempty"`
	SlackChannelId     *string                    `json:"slackChannelId,omitempty"`
	EnrichDomain       *string                    `json:"enrichDomain,omitempty"`
	EnrichSource       *string                    `json:"enrichSource,omitempty"`
	LeadSource         *string                    `json:"leadSource,omitempty"`
	Relationship       *string                    `json:"relationship,omitempty"`
	Stage              *string                    `json:"stage,omitempty"`
	IcpFit             *bool                      `json:"icpFit,omitempty"`
}

func New_UpdateOrganization_From_OrganizationFields(data neo4jrepository.OrganizationSaveFields) UpdateOrganization {
	output := UpdateOrganization{}
	if data.UpdateName {
		output.Name = &data.Name
	}
	if data.UpdateDescription {
		output.Description = &data.Description
	}
	if data.UpdateWebsite {
		output.Website = &data.Website
	}
	if data.UpdateIndustry {
		output.Industry = &data.Industry
	}
	if data.UpdateSubIndustry {
		output.SubIndustry = &data.SubIndustry
	}
	if data.UpdateIndustryGroup {
		output.IndustryGroup = &data.IndustryGroup
	}
	if data.UpdateTargetAudience {
		output.TargetAudience = &data.TargetAudience
	}
	if data.UpdateValueProposition {
		output.ValueProposition = &data.ValueProposition
	}
	if data.UpdateIsPublic {
		output.IsPublic = &data.IsPublic
	}
	if data.UpdateEmployees {
		output.Employees = &data.Employees
	}
	if data.UpdateMarket {
		output.Market = &data.Market
	}
	if data.UpdateLastFundingRound {
		output.LastFundingRound = &data.LastFundingRound
	}
	if data.UpdateLastFundingAmount {
		output.LastFundingAmount = &data.LastFundingAmount
	}
	if data.UpdateCustomerOsId {
		output.CustomerOsId = &data.CustomerOsId
	}
	if data.UpdateReferenceId {
		output.ReferenceId = &data.ReferenceId
	}
	if data.UpdateNote {
		output.Note = &data.Note
	}
	if data.UpdateLogoUrl {
		output.LogoUrl = &data.LogoUrl
	}
	if data.UpdateIconUrl {
		output.IconUrl = &data.IconUrl
	}
	if data.UpdateHeadquarters {
		output.Headquarters = &data.Headquarters
	}
	if data.UpdateYearFounded {
		output.YearFounded = &data.YearFounded
	}
	if data.UpdateEmployeeGrowthRate {
		output.EmployeeGrowthRate = &data.EmployeeGrowthRate
	}
	if data.UpdateSlackChannelId {
		output.SlackChannelId = &data.SlackChannelId
	}
	if data.EnrichDomain != "" {
		output.EnrichDomain = &data.EnrichDomain
	}
	if data.EnrichSource != "" {
		output.EnrichSource = &data.EnrichSource
	}
	if data.UpdateLeadSource {
		output.LeadSource = &data.LeadSource
	}
	if data.UpdateRelationship {
		output.Relationship = utils.StringPtr(data.Relationship.String())
	}
	if data.UpdateStage {
		output.Stage = utils.StringPtr(data.Stage.String())
	}
	if data.UpdateIcpFit {
		output.IcpFit = &data.IcpFit
	}
	if data.ExternalSystem.Available() {
		output.ExternalSystem = &data.ExternalSystem
	}
	return output
}
