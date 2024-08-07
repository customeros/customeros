package opportunity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type OpportunityUpdateNextCycleDateEvent struct {
	Tenant    string     `json:"tenant" validate:"required"`
	UpdatedAt time.Time  `json:"updatedAt"`
	RenewedAt *time.Time `json:"renewedAt"`
}

func NewOpportunityUpdateNextCycleDateEvent(aggregate eventstore.Aggregate, updatedAt time.Time, renewedAt *time.Time) (eventstore.Event, error) {
	eventData := OpportunityUpdateNextCycleDateEvent{
		Tenant:    aggregate.GetTenant(),
		UpdatedAt: updatedAt,
		RenewedAt: utils.ToDatePtr(renewedAt),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OpportunityUpdateNextCycleDateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OpportunityUpdateNextCycleDateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OpportunityUpdateNextCycleDateEvent")
	}
	return event, nil
}
