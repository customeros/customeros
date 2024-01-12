package event

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
)

type OrgPlanMilestoneReorderEvent struct {
	Tenant       string    `json:"tenant" validate:"required"`
	MilestoneIds []string  `json:"milestoneId" validate:"required"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func NewOrgPlanMilestoneReorderEvent(aggregate eventstore.Aggregate, milestoneIds []string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OrgPlanMilestoneReorderEvent{
		Tenant:       aggregate.GetTenant(),
		UpdatedAt:    updatedAt,
		MilestoneIds: milestoneIds,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrgPlanMilestoneReorderEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrgPlanMilestoneReorderV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrgPlanMilestoneReorderEvent")
	}

	return event, nil
}
