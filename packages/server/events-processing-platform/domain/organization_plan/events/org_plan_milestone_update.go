package event

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type OrganizationPlanMilestoneUpdateEvent struct {
	Tenant             string                                `json:"tenant" validate:"required"`
	MilestoneId        string                                `json:"milestoneId" validate:"required"`
	Name               string                                `json:"name,omitempty"`
	Order              int64                                 `json:"order" validate:"gte=0"`
	DueDate            time.Time                             `json:"dueDate"`
	UpdatedAt          time.Time                             `json:"updatedAt"`
	Items              []model.OrganizationPlanMilestoneItem `json:"items"`
	Optional           bool                                  `json:"optional"`
	Retired            bool                                  `json:"retired"`
	FieldsMask         []string                              `json:"fieldsMask,omitempty"`
	OrganizationPlanId string                                `json:"organizationPlanId" validate:"required"`
	StatusDetails      model.OrganizationPlanDetails         `json:"statusDetails"`
	Adhoc              bool                                  `json:"adhoc"`
}

func NewOrganizationPlanMilestoneUpdateEvent(aggregate eventstore.Aggregate, organizationPlanId, milestoneId, name string, order int64, items []model.OrganizationPlanMilestoneItem, fieldsMask []string, optional, adhoc, retired bool, updatedAt, dueDate time.Time, statusDetails model.OrganizationPlanDetails) (eventstore.Event, error) {
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
	if eventData.UpdateDueDate() {
		eventData.DueDate = dueDate
	}
	if eventData.UpdateStatusDetails() {
		eventData.StatusDetails = statusDetails
	}
	if eventData.UpdateAdhoc() {
		eventData.Adhoc = adhoc
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

func (e OrganizationPlanMilestoneUpdateEvent) UpdateDueDate() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskDueDate)
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

func (e OrganizationPlanMilestoneUpdateEvent) UpdateStatusDetails() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskStatusDetails)
}

func (e OrganizationPlanMilestoneUpdateEvent) UpdateAdhoc() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskAdhoc)
}
