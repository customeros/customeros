package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type OrganizationCreateEvent struct {
	Tenant             string                `json:"tenant" validate:"required"`
	Name               string                `json:"name"`
	Hide               bool                  `json:"hide"`
	Description        string                `json:"description"`
	Website            string                `json:"website"`
	Industry           string                `json:"industry"`
	SubIndustry        string                `json:"subIndustry"`
	IndustryGroup      string                `json:"industryGroup"`
	TargetAudience     string                `json:"targetAudience"`
	ValueProposition   string                `json:"valueProposition"`
	IsPublic           bool                  `json:"isPublic"`
	Employees          int64                 `json:"employees"`
	Market             string                `json:"market"`
	LastFundingRound   string                `json:"lastFundingRound"`
	LastFundingAmount  string                `json:"lastFundingAmount"`
	ReferenceId        string                `json:"referenceId"`
	Note               string                `json:"note"`
	Source             string                `json:"source"`
	SourceOfTruth      string                `json:"sourceOfTruth"`
	AppSource          string                `json:"appSource"`
	CreatedAt          time.Time             `json:"createdAt"`
	UpdatedAt          time.Time             `json:"updatedAt"`
	ExternalSystem     cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
	LogoUrl            string                `json:"logoUrl,omitempty"`
	IconUrl            string                `json:"iconUrl,omitempty"`
	YearFounded        *int64                `json:"yearFounded,omitempty"`
	Headquarters       string                `json:"headquarters,omitempty"`
	EmployeeGrowthRate string                `json:"employeeGrowthRate,omitempty"`
	SlackChannelId     string                `json:"slackChannelId,omitempty"`
	Relationship       string                `json:"relationship,omitempty"`
	Stage              string                `json:"stage,omitempty"`
	LeadSource         string                `json:"leadSource,omitempty"`

	// Deprecated
	IsCustomer bool `json:"isCustomer"`
}

func NewOrganizationCreateEvent(aggregate eventstore.Aggregate, organizationFields *model.OrganizationFields, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationCreateEvent{
		Tenant:             aggregate.GetTenant(),
		Name:               organizationFields.OrganizationDataFields.Name,
		Hide:               organizationFields.OrganizationDataFields.Hide,
		Description:        organizationFields.OrganizationDataFields.Description,
		Website:            organizationFields.OrganizationDataFields.Website,
		Industry:           organizationFields.OrganizationDataFields.Industry,
		SubIndustry:        organizationFields.OrganizationDataFields.SubIndustry,
		IndustryGroup:      organizationFields.OrganizationDataFields.IndustryGroup,
		TargetAudience:     organizationFields.OrganizationDataFields.TargetAudience,
		ValueProposition:   organizationFields.OrganizationDataFields.ValueProposition,
		IsPublic:           organizationFields.OrganizationDataFields.IsPublic,
		Employees:          organizationFields.OrganizationDataFields.Employees,
		Market:             organizationFields.OrganizationDataFields.Market,
		LastFundingRound:   organizationFields.OrganizationDataFields.LastFundingRound,
		LastFundingAmount:  organizationFields.OrganizationDataFields.LastFundingAmount,
		ReferenceId:        organizationFields.OrganizationDataFields.ReferenceId,
		Note:               organizationFields.OrganizationDataFields.Note,
		Source:             organizationFields.Source.Source,
		SourceOfTruth:      organizationFields.Source.SourceOfTruth,
		AppSource:          organizationFields.Source.AppSource,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		LogoUrl:            organizationFields.OrganizationDataFields.LogoUrl,
		IconUrl:            organizationFields.OrganizationDataFields.IconUrl,
		YearFounded:        organizationFields.OrganizationDataFields.YearFounded,
		Headquarters:       organizationFields.OrganizationDataFields.Headquarters,
		EmployeeGrowthRate: organizationFields.OrganizationDataFields.EmployeeGrowthRate,
		SlackChannelId:     organizationFields.OrganizationDataFields.SlackChannelId,
		Relationship:       organizationFields.OrganizationDataFields.Relationship,
		Stage:              organizationFields.OrganizationDataFields.Stage,
		LeadSource:         organizationFields.OrganizationDataFields.LeadSource,
	}
	if organizationFields.ExternalSystem.Available() {
		eventData.ExternalSystem = organizationFields.ExternalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationCreateEvent")
	}
	return event, nil
}
