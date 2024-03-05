package event

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type MasterPlanMilestoneCreateEvent struct {
	Tenant        string             `json:"tenant" validate:"required"`
	MilestoneId   string             `json:"milestoneId" validate:"required"`
	Name          string             `json:"name"`
	Order         int64              `json:"order" validate:"gte=0"`
	DurationHours int64              `json:"durationHours" validate:"gte=0"`
	CreatedAt     time.Time          `json:"createdAt"`
	Items         []string           `json:"items"`
	SourceFields  commonmodel.Source `json:"sourceFields"`
	Optional      bool               `json:"optional"`
}

func NewMasterPlanMilestoneCreateEvent(aggregate eventstore.Aggregate, milestoneId, name string, durationHours, order int64, items []string, optional bool, sourceFields commonmodel.Source, createdAt time.Time) (eventstore.Event, error) {
	eventData := MasterPlanMilestoneCreateEvent{
		Tenant:        aggregate.GetTenant(),
		MilestoneId:   milestoneId,
		Name:          name,
		CreatedAt:     createdAt,
		Order:         order,
		DurationHours: durationHours,
		Items:         items,
		SourceFields:  sourceFields,
		Optional:      optional,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate MasterPlanMilestoneCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, MasterPlanMilestoneCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for MasterPlanMilestoneCreateEvent")
	}

	return event, nil
}
