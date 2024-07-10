package event

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type OrganizationPlanMilestoneReorderEvent struct {
	Tenant             string    `json:"tenant" validate:"required"`
	MilestoneIds       []string  `json:"milestoneId" validate:"required"`
	UpdatedAt          time.Time `json:"updatedAt"`
	OrganizationPlanId string    `json:"organizationPlanId" validate:"required"`
}

func NewOrganizationPlanMilestoneReorderEvent(aggregate eventstore.Aggregate, organizationPlanId string, milestoneIds []string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationPlanMilestoneReorderEvent{
		Tenant:             aggregate.GetTenant(),
		UpdatedAt:          updatedAt,
		MilestoneIds:       milestoneIds,
		OrganizationPlanId: organizationPlanId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationPlanMilestoneReorderEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationPlanMilestoneReorderV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationPlanMilestoneReorderEvent")
	}

	return event, nil
}
