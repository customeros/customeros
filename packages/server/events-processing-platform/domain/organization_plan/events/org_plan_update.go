package event

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"gorm.io/gorm/utils"
)

type OrganizationPlanUpdateEvent struct {
	Tenant             string                        `json:"tenant" validate:"required"`
	Name               string                        `json:"name,omitempty"`
	UpdatedAt          time.Time                     `json:"updatedAt"`
	Retired            bool                          `json:"retired"`
	FieldsMask         []string                      `json:"fieldsMask,omitempty"`
	OrganizationPlanId string                        `json:"organizationPlanId" validate:"required"`
	StatusDetails      model.OrganizationPlanDetails `json:"statusDetails"`
}

func NewOrganizationPlanUpdateEvent(aggregate eventstore.Aggregate, organizationPlanId, name string, retired bool, updatedAt time.Time, fieldsMask []string, statusDetails model.OrganizationPlanDetails) (eventstore.Event, error) {
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
	if eventData.UpdateStatusDetails() {
		eventData.StatusDetails = statusDetails
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

func (e OrganizationPlanUpdateEvent) UpdateStatusDetails() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskStatusDetails)
}
