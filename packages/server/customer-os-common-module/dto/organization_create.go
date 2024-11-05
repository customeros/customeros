package dto

import (
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
)

type CreateOrganization struct {
	Source             string                     `json:"source"`
	ExternalSystem     *neo4jmodel.ExternalSystem `json:"externalSystem,omitempty"`
	Domains            []string                   `json:"domains,omitempty"`
	OwnerId            string                     `json:"ownerId"`
	Hide               bool                       `json:"hide"`
	Name               string                     `json:"name"`
	Description        string                     `json:"description"`
	Website            string                     `json:"website"`
	Industry           string                     `json:"industry"`
	SubIndustry        string                     `json:"subIndustry"`
	IndustryGroup      string                     `json:"industryGroup"`
	TargetAudience     string                     `json:"targetAudience"`
	ValueProposition   string                     `json:"valueProposition"`
	IsPublic           bool                       `json:"isPublic"`
	Employees          int64                      `json:"employees"`
	Market             string                     `json:"market"`
	LastFundingRound   string                     `json:"lastFundingRound"`
	LastFundingAmount  string                     `json:"lastFundingAmount"`
	CustomerOsId       string                     `json:"customerOsId"`
	ReferenceId        string                     `json:"referenceId"`
	Note               string                     `json:"note"`
	LogoUrl            string                     `json:"logoUrl"`
	IconUrl            string                     `json:"iconUrl"`
	Headquarters       string                     `json:"headquarters"`
	YearFounded        int64                      `json:"yearFounded"`
	EmployeeGrowthRate string                     `json:"employeeGrowthRate"`
	SlackChannelId     string                     `json:"slackChannelId"`
	EnrichDomain       string                     `json:"enrichDomain"`
	EnrichSource       string                     `json:"enrichSource"`
	LeadSource         string                     `json:"leadSource"`
	Relationship       string                     `json:"relationship"`
	Stage              string                     `json:"stage"`
	IcpFit             bool                       `json:"icpFit"`
}

func New_CreateOrganization_From_OrganizationFields(data neo4jrepository.OrganizationSaveFields) CreateOrganization {
	output := CreateOrganization{
		Source:             data.SourceFields.GetSource(),
		Domains:            data.Domains,
		OwnerId:            data.OwnerId,
		Hide:               data.Hide,
		Name:               data.Name,
		Description:        data.Description,
		Website:            data.Website,
		Industry:           data.Industry,
		SubIndustry:        data.SubIndustry,
		IndustryGroup:      data.IndustryGroup,
		TargetAudience:     data.TargetAudience,
		ValueProposition:   data.ValueProposition,
		IsPublic:           data.IsPublic,
		Employees:          data.Employees,
		Market:             data.Market,
		LastFundingRound:   data.LastFundingRound,
		LastFundingAmount:  data.LastFundingAmount,
		CustomerOsId:       data.CustomerOsId,
		ReferenceId:        data.ReferenceId,
		Note:               data.Note,
		LogoUrl:            data.LogoUrl,
		IconUrl:            data.IconUrl,
		Headquarters:       data.Headquarters,
		YearFounded:        data.YearFounded,
		EmployeeGrowthRate: data.EmployeeGrowthRate,
		SlackChannelId:     data.SlackChannelId,
		EnrichDomain:       data.EnrichDomain,
		EnrichSource:       data.EnrichSource,
		LeadSource:         data.LeadSource,
		Relationship:       data.Relationship.String(),
		Stage:              data.Stage.String(),
		IcpFit:             data.IcpFit,
	}
	if data.ExternalSystem.Available() {
		output.ExternalSystem = &data.ExternalSystem
	}
	return output
}
