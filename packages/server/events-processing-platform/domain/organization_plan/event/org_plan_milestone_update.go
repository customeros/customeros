package event

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
)

type OrgPlanMilestoneUpdateEvent struct {
	Tenant        string                          `json:"tenant" validate:"required"`
	MilestoneId   string                          `json:"milestoneId" validate:"required"`
	Name          string                          `json:"name,omitempty"`
	Order         int64                           `json:"order" validate:"gte=0"`
	DurationHours int64                           `json:"durationHours" validate:"gte=0"`
	UpdatedAt     time.Time                       `json:"updatedAt"`
	Items         []command.OrgPlanMilestoneItems `json:"items"`
	Optional      bool                            `json:"optional"`
	Retired       bool                            `json:"retired"`
	FieldsMask    []string                        `json:"fieldsMask,omitempty"`
}

func NewOrgPlanMilestoneUpdateEvent(aggregate eventstore.Aggregate, milestoneId, name string, durationHours, order int64, items []command.OrgPlanMilestoneItems, fieldsMask []string, optional, retired bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OrgPlanMilestoneUpdateEvent{
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
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrgPlanMilestoneUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrgPlanMilestoneUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrgPlanMilestoneUpdateEvent")
	}

	return event, nil
}

func (e OrgPlanMilestoneUpdateEvent) UpdateName() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskName)
}

func (e OrgPlanMilestoneUpdateEvent) UpdateOrder() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskOrder)
}

func (e OrgPlanMilestoneUpdateEvent) UpdateDurationHours() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskDurationHours)
}

func (e OrgPlanMilestoneUpdateEvent) UpdateItems() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskItems)
}

func (e OrgPlanMilestoneUpdateEvent) UpdateOptional() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskOptional)
}

func (e OrgPlanMilestoneUpdateEvent) UpdateRetired() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskRetired)
}
