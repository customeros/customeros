package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type MasterPlanMilestoneUpdateEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	MilestoneId   string    `json:"milestoneId" validate:"required"`
	Name          string    `json:"name,omitempty"`
	Order         int64     `json:"order" validate:"gte=0"`
	DurationHours int64     `json:"durationHours" validate:"gte=0"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Items         []string  `json:"items"`
	Optional      bool      `json:"optional"`
	Retired       bool      `json:"retired"`
	FieldsMask    []string  `json:"fieldsMask,omitempty"`
}

func NewMasterPlanMilestoneUpdateEvent(aggregate eventstore.Aggregate, milestoneId, name string, durationHours, order int64, items, fieldsMask []string, optional, retired bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := MasterPlanMilestoneUpdateEvent{
		Tenant:      aggregate.GetTenant(),
		MilestoneId: milestoneId,
		UpdatedAt:   updatedAt,
		FieldsMask:  fieldsMask,
	}
	if eventData.UpdateName() {
		eventData.Name = name
	}
	if eventData.UpdateItems() {
		eventData.Items = items
	}
	if eventData.UpdateRetired() {
		eventData.Retired = retired
	}
	if eventData.UpdateOptional() {
		eventData.Optional = optional
	}
	if eventData.UpdateOrder() {
		eventData.Order = order
	}
	if eventData.UpdateDurationHours() {
		eventData.DurationHours = durationHours
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate MasterPlanMilestoneUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, MasterPlanMilestoneUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for MasterPlanMilestoneUpdateEvent")
	}

	return event, nil
}

func (e MasterPlanMilestoneUpdateEvent) UpdateName() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskName)
}

func (e MasterPlanMilestoneUpdateEvent) UpdateOrder() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskOrder)
}

func (e MasterPlanMilestoneUpdateEvent) UpdateDurationHours() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskDurationHours)
}

func (e MasterPlanMilestoneUpdateEvent) UpdateItems() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskItems)
}

func (e MasterPlanMilestoneUpdateEvent) UpdateOptional() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskOptional)
}

func (e MasterPlanMilestoneUpdateEvent) UpdateRetired() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskRetired)
}
