package event

import (
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
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
	RenewalLikelihood string             `json:"renewalLikelihood" validate:"required" enums:"HIGH,MEDIUM,LOW,ZERO"`
}

func NewOpportunityCreateRenewalEvent(aggregate eventstore.Aggregate, contractId, renewalLikelihood string, source commonmodel.Source, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OpportunityCreateRenewalEvent{
		Tenant:            aggregate.GetTenant(),
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		Source:            source,
		ContractId:        contractId,
		InternalType:      neo4jenum.OpportunityInternalTypeRenewal.String(),
		InternalStage:     neo4jenum.OpportunityInternalStageOpen.String(),
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
