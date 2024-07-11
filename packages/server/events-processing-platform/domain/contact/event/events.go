package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

const (
	ContactCreateV1           = "V1_CONTACT_CREATE"
	ContactUpdateV1           = "V1_CONTACT_UPDATE"
	ContactPhoneNumberLinkV1  = "V1_CONTACT_PHONE_NUMBER_LINK"
	ContactEmailLinkV1        = "V1_CONTACT_EMAIL_LINK"
	ContactLocationLinkV1     = "V1_CONTACT_LOCATION_LINK"
	ContactOrganizationLinkV1 = "V1_CONTACT_ORGANIZATION_LINK"
	ContactAddSocialV1        = "V1_CONTACT_ADD_SOCIAL"
	ContactRemoveSocialV1     = "V1_CONTACT_REMOVE_SOCIAL"
	ContactAddTagV1           = "V1_CONTACT_ADD_TAG"
	ContactRemoveTagV1        = "V1_CONTACT_REMOVE_TAG"
	ContactRequestEnrichV1    = "V1_CONTACT_ENRICH"
	ContactAddLocationV1      = "V1_CONTACT_ADD_LOCATION"
	ContactShowV1             = "V1_CONTACT_SHOW"
	ContactHideV1             = "V1_CONTACT_HIDE"
)

const (
	FieldMaskFirstName       = "firstName"
	FieldMaskLastName        = "lastName"
	FieldMaskName            = "name"
	FieldMaskPrefix          = "prefix"
	FieldMaskDescription     = "description"
	FieldMaskTimezone        = "timezone"
	FieldMaskProfilePhotoUrl = "profilePhotoUrl"
)

type ContactCreateEvent struct {
	Tenant          string                `json:"tenant" validate:"required"`
	FirstName       string                `json:"firstName"`
	LastName        string                `json:"lastName"`
	Name            string                `json:"name"`
	Prefix          string                `json:"prefix"`
	Description     string                `json:"description"`
	Timezone        string                `json:"timezone"`
	ProfilePhotoUrl string                `json:"profilePhotoUrl"`
	SocialUrl       string                `json:"socialUrl"`
	Source          string                `json:"source"`
	SourceOfTruth   string                `json:"sourceOfTruth"`
	AppSource       string                `json:"appSource"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	ExternalSystem  cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
}

func NewContactCreateEvent(aggregate eventstore.Aggregate, dataFields models.ContactDataFields, sourceFields events.Source,
	externalSystem cmnmod.ExternalSystem, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactCreateEvent{
		Tenant:          aggregate.GetTenant(),
		FirstName:       dataFields.FirstName,
		LastName:        dataFields.LastName,
		Name:            dataFields.Name,
		Prefix:          dataFields.Prefix,
		Description:     dataFields.Description,
		Timezone:        dataFields.Timezone,
		ProfilePhotoUrl: dataFields.ProfilePhotoUrl,
		SocialUrl:       dataFields.SocialUrl,
		Source:          sourceFields.Source,
		SourceOfTruth:   sourceFields.SourceOfTruth,
		AppSource:       sourceFields.AppSource,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactCreateEvent")
	}
	return event, nil
}

type ContactLinkPhoneNumberEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PhoneNumberId string    `json:"phoneNumberId" validate:"required"`
	Label         string    `json:"label"`
	Primary       bool      `json:"primary"`
}

func NewContactLinkPhoneNumberEvent(aggregate eventstore.Aggregate, phoneNumberId, label string, primary bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactLinkPhoneNumberEvent{
		Tenant:        aggregate.GetTenant(),
		UpdatedAt:     updatedAt,
		PhoneNumberId: phoneNumberId,
		Label:         label,
		Primary:       primary,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactLinkPhoneNumberEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactPhoneNumberLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactLinkPhoneNumberEvent")
	}
	return event, nil
}

type ContactLinkEmailEvent struct {
	Tenant    string    `json:"tenant" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt"`
	EmailId   string    `json:"emailId" validate:"required"`
	Label     string    `json:"label"`
	Primary   bool      `json:"primary"`
}

func NewContactLinkEmailEvent(aggregate eventstore.Aggregate, emailId, label string, primary bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactLinkEmailEvent{
		Tenant:    aggregate.GetTenant(),
		UpdatedAt: updatedAt,
		EmailId:   emailId,
		Label:     label,
		Primary:   primary,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactLinkEmailEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactEmailLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactLinkEmailEvent")
	}
	return event, nil
}

type ContactLinkLocationEvent struct {
	Tenant     string    `json:"tenant" validate:"required"`
	UpdatedAt  time.Time `json:"updatedAt"`
	LocationId string    `json:"locationId" validate:"required"`
}

func NewContactLinkLocationEvent(aggregate eventstore.Aggregate, locationId string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactLinkLocationEvent{
		Tenant:     aggregate.GetTenant(),
		UpdatedAt:  updatedAt,
		LocationId: locationId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactLinkLocationEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactLocationLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactLinkLocationEvent")
	}
	return event, nil
}

type ContactLinkWithOrganizationEvent struct {
	Tenant         string        `json:"tenant" validate:"required"`
	OrganizationId string        `json:"organizationId" validate:"required"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
	StartedAt      *time.Time    `json:"startedAt"`
	EndedAt        *time.Time    `json:"endedAt,omitempty"`
	JobTitle       string        `json:"jobTitle"`
	Description    string        `json:"description"`
	Primary        bool          `json:"primary"`
	SourceFields   events.Source `json:"sourceFields"`
}

func NewContactLinkWithOrganizationEvent(aggregate eventstore.Aggregate, organizationId, jobTile, description string, primary bool,
	sourceFields events.Source, createdAt, updatedAt time.Time, startedAt, endedAt *time.Time) (eventstore.Event, error) {
	eventData := ContactLinkWithOrganizationEvent{
		Tenant:         aggregate.GetTenant(),
		OrganizationId: organizationId,
		SourceFields:   sourceFields,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		StartedAt:      startedAt,
		EndedAt:        endedAt,
		JobTitle:       jobTile,
		Description:    description,
		Primary:        primary,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactLinkWithOrganizationEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactOrganizationLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactLinkWithOrganizationEvent")
	}
	return event, nil
}
