package opportunity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type OpportunityArchiveEvent struct {
	event.BaseEvent
}

func NewOpportunityArchiveEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := OpportunityCloseWinEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate NewOpportunityArchiveEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OpportunityArchiveV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for NewOpportunityArchiveEvent")
	}
	return event, nil
}
