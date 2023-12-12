package events

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
)

type OrganizationUpdateEvent struct {
	// Deprecated
	IgnoreEmptyFields  bool                  `json:"ignoreEmptyFields"`
	Tenant             string                `json:"tenant" validate:"required"`
	Source             string                `json:"source"`
	UpdatedAt          time.Time             `json:"updatedAt"`
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
	IsCustomer         bool                  `json:"isCustomer"`
	Employees          int64                 `json:"employees"`
	Market             string                `json:"market"`
	LastFundingRound   string                `json:"lastFundingRound"`
	LastFundingAmount  string                `json:"lastFundingAmount"`
	ReferenceId        string                `json:"referenceId"`
	Note               string                `json:"note"`
	ExternalSystem     cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
	FieldsMask         []string              `json:"fieldsMask"`
	WebScrapedUrl      string                `json:"webScrapedUrl,omitempty"`
	YearFounded        *int64                `json:"yearFounded"`
	Headquarters       string                `json:"headquarters"`
	EmployeeGrowthRate string                `json:"employeeGrowthRate"`
	LogoUrl            string                `json:"logoUrl"`
}

func NewOrganizationUpdateEvent(aggregate eventstore.Aggregate, organizationFields *model.OrganizationFields, updatedAt time.Time, webScrapedUrl string, fieldsMask []string) (eventstore.Event, error) {
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
		IsCustomer:         organizationFields.OrganizationDataFields.IsCustomer,
		Employees:          organizationFields.OrganizationDataFields.Employees,
		Market:             organizationFields.OrganizationDataFields.Market,
		LastFundingRound:   organizationFields.OrganizationDataFields.LastFundingRound,
		LastFundingAmount:  organizationFields.OrganizationDataFields.LastFundingAmount,
		ReferenceId:        organizationFields.OrganizationDataFields.ReferenceId,
		Note:               organizationFields.OrganizationDataFields.Note,
		LogoUrl:            organizationFields.OrganizationDataFields.LogoUrl,
		YearFounded:        organizationFields.OrganizationDataFields.YearFounded,
		Headquarters:       organizationFields.OrganizationDataFields.Headquarters,
		EmployeeGrowthRate: organizationFields.OrganizationDataFields.EmployeeGrowthRate,
		UpdatedAt:          updatedAt,
		Source:             organizationFields.Source.Source,
		FieldsMask:         fieldsMask,
		WebScrapedUrl:      webScrapedUrl,
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

func (e OrganizationUpdateEvent) shouldUpdateFieldIfNotIgnored(input string) bool {
	return e.IgnoreEmptyFields == false || input != ""
}

func (e OrganizationUpdateEvent) UpdateName() bool {
	return (len(e.FieldsMask) == 0 && e.shouldUpdateFieldIfNotIgnored(e.Name)) || utils.Contains(e.FieldsMask, model.FieldMaskName)
}

func (e OrganizationUpdateEvent) UpdateHide() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskHide)
}

func (e OrganizationUpdateEvent) UpdateDescription() bool {
	return (len(e.FieldsMask) == 0 && e.shouldUpdateFieldIfNotIgnored(e.Description)) || utils.Contains(e.FieldsMask, model.FieldMaskDescription)
}

func (e OrganizationUpdateEvent) UpdateWebsite() bool {
	return (len(e.FieldsMask) == 0 && e.shouldUpdateFieldIfNotIgnored(e.Website)) || utils.Contains(e.FieldsMask, model.FieldMaskWebsite)
}

func (e OrganizationUpdateEvent) UpdateIndustry() bool {
	return (len(e.FieldsMask) == 0 && e.shouldUpdateFieldIfNotIgnored(e.Industry)) || utils.Contains(e.FieldsMask, model.FieldMaskIndustry)
}

func (e OrganizationUpdateEvent) UpdateSubIndustry() bool {
	return (len(e.FieldsMask) == 0 && e.shouldUpdateFieldIfNotIgnored(e.SubIndustry)) || utils.Contains(e.FieldsMask, model.FieldMaskSubIndustry)
}

func (e OrganizationUpdateEvent) UpdateIndustryGroup() bool {
	return (len(e.FieldsMask) == 0 && e.shouldUpdateFieldIfNotIgnored(e.IndustryGroup)) || utils.Contains(e.FieldsMask, model.FieldMaskIndustryGroup)
}

func (e OrganizationUpdateEvent) UpdateTargetAudience() bool {
	return (len(e.FieldsMask) == 0 && e.shouldUpdateFieldIfNotIgnored(e.TargetAudience)) || utils.Contains(e.FieldsMask, model.FieldMaskTargetAudience)
}

func (e OrganizationUpdateEvent) UpdateValueProposition() bool {
	return (len(e.FieldsMask) == 0 && e.shouldUpdateFieldIfNotIgnored(e.ValueProposition)) || utils.Contains(e.FieldsMask, model.FieldMaskValueProposition)
}

func (e OrganizationUpdateEvent) UpdateIsPublic() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskIsPublic)
}

func (e OrganizationUpdateEvent) UpdateIsCustomer() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskIsCustomer)
}

func (e OrganizationUpdateEvent) UpdateEmployees() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskEmployees)
}

func (e OrganizationUpdateEvent) UpdateMarket() bool {
	return (len(e.FieldsMask) == 0 && e.shouldUpdateFieldIfNotIgnored(e.Market)) || utils.Contains(e.FieldsMask, model.FieldMaskMarket)
}

func (e OrganizationUpdateEvent) UpdateLastFundingRound() bool {
	return (len(e.FieldsMask) == 0 && e.shouldUpdateFieldIfNotIgnored(e.LastFundingRound)) || utils.Contains(e.FieldsMask, model.FieldMaskLastFundingRound)
}

func (e OrganizationUpdateEvent) UpdateLastFundingAmount() bool {
	return (len(e.FieldsMask) == 0 && e.shouldUpdateFieldIfNotIgnored(e.LastFundingAmount)) || utils.Contains(e.FieldsMask, model.FieldMaskLastFundingAmount)
}

func (e OrganizationUpdateEvent) UpdateReferenceId() bool {
	return (len(e.FieldsMask) == 0 && e.shouldUpdateFieldIfNotIgnored(e.ReferenceId)) || utils.Contains(e.FieldsMask, model.FieldMaskReferenceId)
}

func (e OrganizationUpdateEvent) UpdateNote() bool {
	return (len(e.FieldsMask) == 0 && e.shouldUpdateFieldIfNotIgnored(e.Note)) || utils.Contains(e.FieldsMask, model.FieldMaskNote)
}

func (e OrganizationUpdateEvent) UpdateYearFounded() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskYearFounded)
}

func (e OrganizationUpdateEvent) UpdateHeadquarters() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskHeadquarters)
}

func (e OrganizationUpdateEvent) UpdateEmployeeGrowthRate() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskEmployeeGrowthRate)
}

func (e OrganizationUpdateEvent) UpdateLogoUrl() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskLogoUrl)
}
