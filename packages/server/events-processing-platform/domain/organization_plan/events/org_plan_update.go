package event

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"gorm.io/gorm/utils"
)

type OrganizationPlanUpdateEvent struct {
	Tenant             string    `json:"tenant" validate:"required"`
	Name               string    `json:"name,omitempty"`
	UpdatedAt          time.Time `json:"updatedAt"`
	Retired            bool      `json:"retired"`
	FieldsMask         []string  `json:"fieldsMask,omitempty"`
	OrganizationPlanId string    `json:"organizationPlanId" validate:"required"`
}

func NewOrganizationPlanUpdateEvent(aggregate eventstore.Aggregate, organizationPlanId, name string, retired bool, updatedAt time.Time, fieldsMask []string) (eventstore.Event, error) {
	eventData := OrganizationPlanUpdateEvent{
		Tenant:             aggregate.GetTenant(),
		UpdatedAt:          updatedAt,
		FieldsMask:         fieldsMask,
		OrganizationPlanId: organizationPlanId,
	}
	if eventData.UpdateName() {
		eventData.Name = name
	}
	if eventData.UpdateRetired() {
		eventData.Retired = retired
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationPlanUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationPlanUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationPlanUpdateEvent")
	}

	return event, nil
}

func (e OrganizationPlanUpdateEvent) UpdateName() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskName)
}

func (e OrganizationPlanUpdateEvent) UpdateRetired() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskRetired)
}
