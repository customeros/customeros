package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type ContactUpdateEvent struct {
	Tenant          string                `json:"tenant" validate:"required"`
	Source          string                `json:"source"`
	AppSource       string                `json:"appSource"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	FirstName       string                `json:"firstName"`
	LastName        string                `json:"lastName"`
	Name            string                `json:"name"`
	Prefix          string                `json:"prefix"`
	Description     string                `json:"description"`
	Timezone        string                `json:"timezone"`
	ProfilePhotoUrl string                `json:"profilePhotoUrl"`
	Username        string                `json:"username"`
	ExternalSystem  cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
	FieldsMask      []string              `json:"fieldsMask,omitempty"`
}

func NewContactUpdateEvent(aggregate eventstore.Aggregate, source, appSource string, dataFields ContactDataFields, externalSystem cmnmod.ExternalSystem, updatedAt time.Time, fieldsMask []string) (eventstore.Event, error) {
	eventData := ContactUpdateEvent{
		Tenant:          aggregate.GetTenant(),
		FirstName:       dataFields.FirstName,
		LastName:        dataFields.LastName,
		Prefix:          dataFields.Prefix,
		Description:     dataFields.Description,
		Timezone:        dataFields.Timezone,
		ProfilePhotoUrl: dataFields.ProfilePhotoUrl,
		Username:        dataFields.Username,
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

func (e ContactUpdateEvent) UpdateUsername() bool {
	return utils.Contains(e.FieldsMask, FieldMaskUsername)
}
