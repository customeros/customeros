package event

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type OpportunityCreateRenewalEvent struct {
	Tenant            string             `json:"tenant" validate:"required"`
	CreatedAt         time.Time          `json:"createdAt"`
	UpdatedAt         time.Time          `json:"updatedAt"`
	Source            commonmodel.Source `json:"source"`
	ContractId        string             `json:"contractId" validate:"required"`
	InternalType      string             `json:"internalType"`
	InternalStage     string             `json:"internalStage"`
	RenewalLikelihood string             `json:"renewalLikelihood"`
}

func NewOpportunityCreateRenewalEvent(aggregate eventstore.Aggregate, contractId, renewalLikelihood string, source commonmodel.Source, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OpportunityCreateRenewalEvent{
		Tenant:            aggregate.GetTenant(),
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		Source:            source,
		ContractId:        contractId,
		InternalType:      string(model.OpportunityInternalTypeStringRenewal),
		InternalStage:     string(model.OpportunityInternalStageStringOpen),
		RenewalLikelihood: renewalLikelihood,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OpportunityCreateRenewalEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OpportunityCreateRenewalV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OpportunityCreateRenewalEvent")
	}
	return event, nil
}
