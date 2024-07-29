package events

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type OrganizationUpdateEvent struct {
	Tenant             string                `json:"tenant" validate:"required"`
	Source             string                `json:"source,omitempty"`
	UpdatedAt          time.Time             `json:"updatedAt,omitempty"`
	Name               string                `json:"name,omitempty"`
	Hide               bool                  `json:"hide,omitempty"`
	Description        string                `json:"description,omitempty"`
	Website            string                `json:"website,omitempty"`
	Industry           string                `json:"industry,omitempty"`
	SubIndustry        string                `json:"subIndustry,omitempty"`
	IndustryGroup      string                `json:"industryGroup,omitempty"`
	TargetAudience     string                `json:"targetAudience,omitempty"`
	ValueProposition   string                `json:"valueProposition,omitempty"`
	IsPublic           bool                  `json:"isPublic,omitempty"`
	Employees          int64                 `json:"employees,omitempty"`
	Market             string                `json:"market,omitempty"`
	LastFundingRound   string                `json:"lastFundingRound,omitempty"`
	LastFundingAmount  string                `json:"lastFundingAmount,omitempty"`
	ReferenceId        string                `json:"referenceId,omitempty"`
	Note               string                `json:"note,omitempty"`
	ExternalSystem     cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
	FieldsMask         []string              `json:"fieldsMask"`
	YearFounded        *int64                `json:"yearFounded,omitempty"`
	Headquarters       string                `json:"headquarters,omitempty"`
	EmployeeGrowthRate string                `json:"employeeGrowthRate,omitempty"`
	SlackChannelId     string                `json:"slackChannelId,omitempty"`
	LogoUrl            string                `json:"logoUrl,omitempty"`
	IconUrl            string                `json:"iconUrl,omitempty"`
	Relationship       string                `json:"relationship,omitempty"`
	Stage              string                `json:"stage,omitempty"`
	EnrichDomain       string                `json:"enrichDomain,omitempty"`
	EnrichSource       string                `json:"enrichSource,omitempty"`

	// Deprecated
	IsCustomer bool `json:"isCustomer,omitempty"`
	// Deprecated
	IgnoreEmptyFields bool `json:"ignoreEmptyFields"`
	// Deprecated
	WebScrapedUrl string `json:"webScrapedUrl,omitempty"`
}

func NewOrganizationUpdateEvent(aggregate eventstore.Aggregate, organizationFields *model.OrganizationFields, updatedAt time.Time, enrichDomain, enrichSource string, fieldsMask []string) (eventstore.Event, error) {
	eventData := OrganizationUpdateEvent{
		IgnoreEmptyFields:  false,
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
		LogoUrl:            organizationFields.OrganizationDataFields.LogoUrl,
		IconUrl:            organizationFields.OrganizationDataFields.IconUrl,
		YearFounded:        organizationFields.OrganizationDataFields.YearFounded,
		Headquarters:       organizationFields.OrganizationDataFields.Headquarters,
		EmployeeGrowthRate: organizationFields.OrganizationDataFields.EmployeeGrowthRate,
		SlackChannelId:     organizationFields.OrganizationDataFields.SlackChannelId,
		UpdatedAt:          updatedAt,
		Source:             organizationFields.Source.Source,
		FieldsMask:         fieldsMask,
		EnrichDomain:       enrichDomain,
		EnrichSource:       enrichSource,
		Relationship:       organizationFields.OrganizationDataFields.Relationship,
		Stage:              organizationFields.OrganizationDataFields.Stage,
	}
	if organizationFields.ExternalSystem.Available() {
		eventData.ExternalSystem = organizationFields.ExternalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationUpdateEvent")
	}
	return event, nil
}

func (e OrganizationUpdateEvent) UpdateName() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskName)
}

func (e OrganizationUpdateEvent) UpdateHide() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskHide)
}

func (e OrganizationUpdateEvent) UpdateDescription() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskDescription)
}

func (e OrganizationUpdateEvent) UpdateWebsite() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskWebsite)
}

func (e OrganizationUpdateEvent) UpdateIndustry() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskIndustry)
}

func (e OrganizationUpdateEvent) UpdateSubIndustry() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskSubIndustry)
}

func (e OrganizationUpdateEvent) UpdateIndustryGroup() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskIndustryGroup)
}

func (e OrganizationUpdateEvent) UpdateTargetAudience() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskTargetAudience)
}

func (e OrganizationUpdateEvent) UpdateValueProposition() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskValueProposition)
}

func (e OrganizationUpdateEvent) UpdateIsPublic() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskIsPublic)
}

func (e OrganizationUpdateEvent) UpdateEmployees() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskEmployees)
}

func (e OrganizationUpdateEvent) UpdateMarket() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskMarket)
}

func (e OrganizationUpdateEvent) UpdateLastFundingRound() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskLastFundingRound)
}

func (e OrganizationUpdateEvent) UpdateLastFundingAmount() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskLastFundingAmount)
}

func (e OrganizationUpdateEvent) UpdateReferenceId() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskReferenceId)
}

func (e OrganizationUpdateEvent) UpdateNote() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskNote)
}

func (e OrganizationUpdateEvent) UpdateYearFounded() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskYearFounded)
}

func (e OrganizationUpdateEvent) UpdateHeadquarters() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskHeadquarters)
}

func (e OrganizationUpdateEvent) UpdateEmployeeGrowthRate() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskEmployeeGrowthRate)
}

func (e OrganizationUpdateEvent) UpdateSlackChannelId() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskSlackChannelId)
}

func (e OrganizationUpdateEvent) UpdateLogoUrl() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskLogoUrl)
}

func (e OrganizationUpdateEvent) UpdateIconUrl() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskIconUrl)
}

func (e OrganizationUpdateEvent) UpdateRelationship() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskRelationship)
}

func (e OrganizationUpdateEvent) UpdateStage() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskStage)
}
