package event

import (
	"time"

	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
)

type OrgPlanCreateEvent struct {
	Tenant       string             `json:"tenant" validate:"required"`
	Name         string             `json:"name"`
	CreatedAt    time.Time          `json:"createdAt"`
	SourceFields commonmodel.Source `json:"sourceFields"`
}

func NewOrgPlanCreateEvent(aggregate eventstore.Aggregate, name string, sourceFields commonmodel.Source, createdAt time.Time) (eventstore.Event, error) {
	eventData := OrgPlanCreateEvent{
		Tenant:       aggregate.GetTenant(),
		Name:         name,
		CreatedAt:    createdAt,
		SourceFields: sourceFields,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrgPlanCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrgPlanCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrgPlanCreateEvent")
	}

	return event, nil
}
