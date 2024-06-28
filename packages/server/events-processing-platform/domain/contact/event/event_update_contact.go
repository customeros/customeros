package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
	"time"
)

type ContactUpdateEvent struct {
	Tenant          string                `json:"tenant" validate:"required"`
	Source          string                `json:"source"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	FirstName       string                `json:"firstName"`
	LastName        string                `json:"lastName"`
	Name            string                `json:"name"`
	Prefix          string                `json:"prefix"`
	Description     string                `json:"description"`
	Timezone        string                `json:"timezone"`
	ProfilePhotoUrl string                `json:"profilePhotoUrl"`
	ExternalSystem  cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
	FieldsMask      []string              `json:"fieldsMask,omitempty"`
}

func NewContactUpdateEvent(aggregate eventstore.Aggregate, source string, dataFields models.ContactDataFields, externalSystem cmnmod.ExternalSystem, updatedAt time.Time, fieldsMask []string) (eventstore.Event, error) {
	eventData := ContactUpdateEvent{
		Tenant:          aggregate.GetTenant(),
		FirstName:       dataFields.FirstName,
		LastName:        dataFields.LastName,
		Prefix:          dataFields.Prefix,
		Description:     dataFields.Description,
		Timezone:        dataFields.Timezone,
		ProfilePhotoUrl: dataFields.ProfilePhotoUrl,
		Name:            dataFields.Name,
		UpdatedAt:       updatedAt,
		Source:          source,
		FieldsMask:      fieldsMask,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactUpdateEvent")
	}
	return event, nil
}

func (e ContactUpdateEvent) UpdateFirstName() bool {
	return utils.Contains(e.FieldsMask, FieldMaskFirstName)
}

func (e ContactUpdateEvent) UpdateLastName() bool {
	return utils.Contains(e.FieldsMask, FieldMaskLastName)
}

func (e ContactUpdateEvent) UpdateName() bool {
	return utils.Contains(e.FieldsMask, FieldMaskName)
}

func (e ContactUpdateEvent) UpdatePrefix() bool {
	return utils.Contains(e.FieldsMask, FieldMaskPrefix)
}

func (e ContactUpdateEvent) UpdateDescription() bool {
	return utils.Contains(e.FieldsMask, FieldMaskDescription)
}

func (e ContactUpdateEvent) UpdateTimezone() bool {
	return utils.Contains(e.FieldsMask, FieldMaskTimezone)
}

func (e ContactUpdateEvent) UpdateProfilePhotoUrl() bool {
	return utils.Contains(e.FieldsMask, FieldMaskProfilePhotoUrl)
}
