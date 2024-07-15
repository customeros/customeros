package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type OpportunityCreateRenewalEvent struct {
	Tenant              string        `json:"tenant" validate:"required"`
	CreatedAt           time.Time     `json:"createdAt"`
	UpdatedAt           time.Time     `json:"updatedAt"`
	Source              common.Source `json:"source"`
	ContractId          string        `json:"contractId" validate:"required"`
	InternalType        string        `json:"internalType"`
	InternalStage       string        `json:"internalStage"`
	RenewalLikelihood   string        `json:"renewalLikelihood" validate:"required" enums:"HIGH,MEDIUM,LOW,ZERO"`
	RenewalApproved     bool          `json:"renewalApproved,omitempty"`
	RenewedAt           *time.Time    `json:"renewedAt,omitempty"`
	RenewalAdjustedRate int64         `json:"renewalAdjustedRate" validate:"min=0,max=100"`
}

func NewOpportunityCreateRenewalEvent(aggregate eventstore.Aggregate, contractId, renewalLikelihood string, renewalApproved bool, source common.Source, createdAt, updatedAt time.Time, renewedAt *time.Time, adjustedRate int64) (eventstore.Event, error) {
	eventData := OpportunityCreateRenewalEvent{
		Tenant:              aggregate.GetTenant(),
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
		Source:              source,
		ContractId:          contractId,
		InternalType:        neo4jenum.OpportunityInternalTypeRenewal.String(),
		InternalStage:       neo4jenum.OpportunityInternalStageOpen.String(),
		RenewalLikelihood:   renewalLikelihood,
		RenewalApproved:     renewalApproved,
		RenewedAt:           renewedAt,
		RenewalAdjustedRate: adjustedRate,
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
