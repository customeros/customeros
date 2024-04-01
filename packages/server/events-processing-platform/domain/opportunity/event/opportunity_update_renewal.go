package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
	"time"
)

type OpportunityUpdateRenewalEvent struct {
	Tenant            string    `json:"tenant" validate:"required"`
	UpdatedAt         time.Time `json:"updatedAt"`
	UpdatedByUserId   string    `json:"updatedByUserId"`
	RenewalLikelihood string    `json:"renewalLikelihood" validate:"required" enums:"HIGH,MEDIUM,LOW,ZERO"`
	RenewalApproved   bool      `json:"renewalApproved,omitempty"`
	Comments          string    `json:"comments"`
	Amount            float64   `json:"amount"`
	Source            string    `json:"source"`
	FieldsMask        []string  `json:"fieldsMask"`
	OwnerUserId       string    `json:"ownerUserId"`
}

func NewOpportunityUpdateRenewalEvent(aggregate eventstore.Aggregate, renewalLikelihood, comments, updatedByUserId, source string, amount float64, renewalApproved bool, updatedAt time.Time, fieldsMask []string, ownerUserId string) (eventstore.Event, error) {
	eventData := OpportunityUpdateRenewalEvent{
		Tenant:            aggregate.GetTenant(),
		UpdatedAt:         updatedAt,
		Source:            source,
		RenewalLikelihood: renewalLikelihood,
		RenewalApproved:   renewalApproved,
		Comments:          comments,
		UpdatedByUserId:   updatedByUserId,
		Amount:            amount,
		FieldsMask:        fieldsMask,
		OwnerUserId:       ownerUserId,
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

func (e OpportunityUpdateRenewalEvent) UpdateAmount() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskAmount)
}

func (e OpportunityUpdateRenewalEvent) UpdateComments() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskComments)
}

func (e OpportunityUpdateRenewalEvent) UpdateRenewalLikelihood() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskRenewalLikelihood)
}

func (e OpportunityUpdateRenewalEvent) UpdateRenewalApproved() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskRenewalApproved)
}
