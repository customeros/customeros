package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type OpportunityUpdateRenewalEvent struct {
	Tenant            string    `json:"tenant" validate:"required"`
	UpdatedAt         time.Time `json:"updatedAt"`
	UpdatedByUserId   string    `json:"updatedByUserId"`
	RenewalLikelihood string    `json:"renewalLikelihood" validate:"required" enums:"HIGH,MEDIUM,LOW,ZERO"`
	Comments          string    `json:"comments"`
	Amount            float64   `json:"amount"`
	Source            string    `json:"source"`
}

func NewOpportunityUpdateRenewalEvent(aggregate eventstore.Aggregate, renewalLikelihood, comments, updatedByUserId, source string, amount float64, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OpportunityUpdateRenewalEvent{
		Tenant:            aggregate.GetTenant(),
		UpdatedAt:         updatedAt,
		Source:            source,
		RenewalLikelihood: renewalLikelihood,
		Comments:          comments,
		UpdatedByUserId:   updatedByUserId,
		Amount:            amount,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OpportunityUpdateRenewalEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OpportunityUpdateRenewalV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OpportunityUpdateRenewalEvent")
	}
	return event, nil
}
