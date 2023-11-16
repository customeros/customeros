package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type OpportunityUpdateEvent struct {
	Tenant         string                     `json:"tenant" validate:"required"`
	Name           string                     `json:"name"`
	Amount         float64                    `json:"amount"`
	MaxAmount      float64                    `json:"maxAmount"`
	UpdatedAt      time.Time                  `json:"updatedAt"`
	Source         string                     `json:"source"`
	ExternalSystem commonmodel.ExternalSystem `json:"externalSystem,omitempty"`
	FieldsMask     []string                   `json:"fieldsMask"`
}

func NewOpportunityUpdateEvent(aggregate eventstore.Aggregate, dataFields model.OpportunityDataFields, source string, externalSystem commonmodel.ExternalSystem, updatedAt time.Time, fieldsMask []string) (eventstore.Event, error) {
	eventData := OpportunityUpdateEvent{
		Tenant:     aggregate.GetTenant(),
		Name:       dataFields.Name,
		Amount:     dataFields.Amount,
		MaxAmount:  dataFields.MaxAmount,
		UpdatedAt:  updatedAt,
		Source:     source,
		FieldsMask: fieldsMask,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OpportunityUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OpportunityUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OpportunityUpdateEvent")
	}
	return event, nil
}

func (e OpportunityUpdateEvent) UpdateName() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskName)
}

func (e OpportunityUpdateEvent) UpdateAmount() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskAmount)
}

func (e OpportunityUpdateEvent) UpdateMaxAmount() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, model.FieldMaskMaxAmount)
}
