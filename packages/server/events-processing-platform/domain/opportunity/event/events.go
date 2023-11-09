package event

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

const (
	OpportunityCreateV1 = "V1_OPPORTUNITY_CREATE"
	// OpportunityUpdateV1 = "V1_OPPORTUNITY_UPDATE" // Define this for update events.
)

type OpportunityCreateEvent struct {
	Tenant            string                         `json:"tenant" validate:"required"`
	Name              string                         `json:"name"`
	Amount            float64                        `json:"amount"`
	InternalType      model.OpportunityInternalType  `json:"internalType"`
	ExternalType      string                         `json:"externalType"`
	InternalStage     model.OpportunityInternalStage `json:"internalStage"`
	ExternalStage     string                         `json:"externalStage"`
	EstimatedClosedAt *time.Time                     `json:"estimatedClosedAt,omitempty"`
	OwnerUserId       string                         `json:"ownerUserId"`
	CreatedByUserId   string                         `json:"createdByUserId"`
	CreatedAt         time.Time                      `json:"createdAt"`
	UpdatedAt         time.Time                      `json:"updatedAt"`
	Source            commonmodel.Source             `json:"source"`
	ExternalSystem    commonmodel.ExternalSystem     `json:"externalSystem,omitempty"`
	OrganizationId    string                         `json:"organizationId" validate:"required"`
	GeneralNotes      string                         `json:"generalNotes"`
	NextSteps         string                         `json:"nextSteps"`
}

func NewOpportunityCreateEvent(aggregate eventstore.Aggregate, dataFields model.OpportunityDataFields, source commonmodel.Source, externalSystem commonmodel.ExternalSystem, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OpportunityCreateEvent{
		Tenant:            aggregate.GetTenant(),
		Name:              dataFields.Name,
		Amount:            dataFields.Amount,
		InternalType:      dataFields.InternalType.StrValue(),
		ExternalType:      dataFields.ExternalType,
		InternalStage:     dataFields.InternalStage.StrValue(),
		ExternalStage:     dataFields.ExternalStage,
		EstimatedClosedAt: dataFields.EstimatedClosedAt,
		OwnerUserId:       dataFields.OwnerUserId,
		CreatedByUserId:   dataFields.CreatedByUserId,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		Source:            source,
		OrganizationId:    dataFields.OrganizationId,
		GeneralNotes:      dataFields.GeneralNotes,
		NextSteps:         dataFields.NextSteps,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OpportunityCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OpportunityCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OpportunityCreateEvent")
	}
	return event, nil
}
