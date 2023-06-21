package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"time"
)

const (
	OrganizationCreateV1          = "V1_ORGANIZATION_CREATE"
	OrganizationUpdateV1          = "V1_ORGANIZATION_UPDATE"
	OrganizationPhoneNumberLinkV1 = "V1_ORGANIZATION_PHONE_NUMBER_LINK"
	OrganizationEmailLinkV1       = "V1_ORGANIZATION_EMAIL_LINK"
)

type OrganizationCreateEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	Name          string    `json:"name" required:"true"`
	Description   string    `json:"description"`
	Website       string    `json:"website"`
	Industry      string    `json:"industry"`
	IsPublic      bool      `json:"isPublic"`
	Employees     int64     `json:"employees"`
	Market        string    `json:"market"`
	Source        string    `json:"source"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	AppSource     string    `json:"appSource"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewOrganizationCreateEvent(aggregate eventstore.Aggregate, organizationDto *models.OrganizationDto, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationCreateEvent{
		Tenant:        organizationDto.Tenant,
		Name:          organizationDto.Name,
		Description:   organizationDto.Description,
		Website:       organizationDto.Website,
		Industry:      organizationDto.Industry,
		IsPublic:      organizationDto.IsPublic,
		Employees:     organizationDto.Employees,
		Market:        organizationDto.Market,
		Source:        organizationDto.Source.Source,
		SourceOfTruth: organizationDto.Source.SourceOfTruth,
		AppSource:     organizationDto.Source.AppSource,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type OrganizationUpdateEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Website       string    `json:"website"`
	Industry      string    `json:"industry"`
	IsPublic      bool      `json:"isPublic"`
	Employees     int64     `json:"employees"`
	Market        string    `json:"market"`
}

func NewOrganizationUpdateEvent(aggregate eventstore.Aggregate, organizationDto *models.OrganizationDto, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationUpdateEvent{
		Name:          organizationDto.Name,
		Description:   organizationDto.Description,
		Website:       organizationDto.Website,
		Industry:      organizationDto.Industry,
		IsPublic:      organizationDto.IsPublic,
		Tenant:        organizationDto.Tenant,
		Employees:     organizationDto.Employees,
		Market:        organizationDto.Market,
		UpdatedAt:     updatedAt,
		SourceOfTruth: organizationDto.Source.SourceOfTruth,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type OrganizationLinkPhoneNumberEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PhoneNumberId string    `json:"phoneNumberId" validate:"required"`
	Label         string    `json:"label"`
	Primary       bool      `json:"primary"`
}

func NewOrganizationLinkPhoneNumberEvent(aggregate eventstore.Aggregate, tenant, phoneNumberId, label string, primary bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationLinkPhoneNumberEvent{
		Tenant:        tenant,
		UpdatedAt:     updatedAt,
		PhoneNumberId: phoneNumberId,
		Label:         label,
		Primary:       primary,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationPhoneNumberLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type OrganizationLinkEmailEvent struct {
	Tenant    string    `json:"tenant" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt"`
	EmailId   string    `json:"emailId" validate:"required"`
	Label     string    `json:"label"`
	Primary   bool      `json:"primary"`
}

func NewOrganizationLinkEmailEvent(aggregate eventstore.Aggregate, tenant, emailId, label string, primary bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationLinkEmailEvent{
		Tenant:    tenant,
		UpdatedAt: updatedAt,
		EmailId:   emailId,
		Label:     label,
		Primary:   primary,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationEmailLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
