package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
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
	ExternalStage  string                     `json:"externalStage"`
	ExternalType   string                     `json:"externalType"`
}

func NewOpportunityUpdateEvent(aggregate eventstore.Aggregate, dataFields model.OpportunityDataFields, source string, externalSystem commonmodel.ExternalSystem, updatedAt time.Time, fieldsMask []string) (eventstore.Event, error) {
	eventData := OpportunityUpdateEvent{
		Tenant:        aggregate.GetTenant(),
		Name:          dataFields.Name,
		Amount:        dataFields.Amount,
		MaxAmount:     dataFields.MaxAmount,
		ExternalStage: dataFields.ExternalStage,
		ExternalType:  dataFields.ExternalType,
		UpdatedAt:     updatedAt,
		Source:        source,
		FieldsMask:    fieldsMask,
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
	return utils.Contains(e.FieldsMask, model.FieldMaskName)
}

func (e OpportunityUpdateEvent) UpdateAmount() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskAmount)
}

func (e OpportunityUpdateEvent) UpdateMaxAmount() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskMaxAmount)
}

func (e OpportunityUpdateEvent) UpdateExternalStage() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskExternalStage)
}

func (e OpportunityUpdateEvent) UpdateExternalType() bool {
	return utils.Contains(e.FieldsMask, model.FieldMaskExternalType)
}
