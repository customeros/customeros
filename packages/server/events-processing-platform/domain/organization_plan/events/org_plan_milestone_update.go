package event

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
)

type OrganizationPlanMilestoneUpdateEvent struct {
	Tenant             string                                `json:"tenant" validate:"required"`
	MilestoneId        string                                `json:"milestoneId" validate:"required"`
	Name               string                                `json:"name,omitempty"`
	Order              int64                                 `json:"order" validate:"gte=0"`
	DurationHours      int64                                 `json:"durationHours" validate:"gte=0"`
	UpdatedAt          time.Time                             `json:"updatedAt"`
	Items              []model.OrganizationPlanMilestoneItem `json:"items"`
	Optional           bool                                  `json:"optional"`
	Retired            bool                                  `json:"retired"`
	FieldsMask         []string                              `json:"fieldsMask,omitempty"`
	OrganizationPlanId string                                `json:"organizationPlanId" validate:"required"`
}

func NewOrganizationPlanMilestoneUpdateEvent(aggregate eventstore.Aggregate, organizationPlanId, milestoneId, name string, durationHours, order int64, items []model.OrganizationPlanMilestoneItem, fieldsMask []string, optional, retired bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationPlanMilestoneUpdateEvent{
		Tenant:             aggregate.GetTenant(),
		MilestoneId:        milestoneId,
		UpdatedAt:          updatedAt,
		FieldsMask:         fieldsMask,
		OrganizationPlanId: organizationPlanId,
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
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationPlanMilestoneUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationPlanMilestoneUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationPlanMilestoneUpdateEvent")
	}

	return event, nil
}

func (e OrganizationPlanMilestoneUpdateEvent) UpdateName() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskName)
}

func (e OrganizationPlanMilestoneUpdateEvent) UpdateOrder() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskOrder)
}

func (e OrganizationPlanMilestoneUpdateEvent) UpdateDurationHours() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskDurationHours)
}

func (e OrganizationPlanMilestoneUpdateEvent) UpdateItems() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskItems)
}

func (e OrganizationPlanMilestoneUpdateEvent) UpdateOptional() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskOptional)
}

func (e OrganizationPlanMilestoneUpdateEvent) UpdateRetired() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskRetired)
}
