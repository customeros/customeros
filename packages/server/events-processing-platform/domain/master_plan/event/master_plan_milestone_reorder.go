package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type MasterPlanMilestoneReorderEvent struct {
	Tenant       string    `json:"tenant" validate:"required"`
	MilestoneIds []string  `json:"milestoneId" validate:"required"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func NewMasterPlanMilestoneReorderEvent(aggregate eventstore.Aggregate, milestoneIds []string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := MasterPlanMilestoneReorderEvent{
		Tenant:       aggregate.GetTenant(),
		UpdatedAt:    updatedAt,
		MilestoneIds: milestoneIds,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate MasterPlanMilestoneReorderEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, MasterPlanMilestoneReorderV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for MasterPlanMilestoneReorderEvent")
	}

	return event, nil
}
