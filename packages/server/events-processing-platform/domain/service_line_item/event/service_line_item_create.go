package event

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type ServiceLineItemCreateEvent struct {
	Tenant     string             `json:"tenant" validate:"required"`
	Billed     string             `json:"billed"`
	Quantity   int64              `json:"quantity,omitempty"`
	Price      float64            `json:"price"`
	Name       string             `json:"name"`
	ContractId string             `json:"contractId" validate:"required"`
	ParentId   string             `json:"parentId" validate:"required"`
	CreatedAt  time.Time          `json:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt"`
	StartedAt  time.Time          `json:"startedAt"`
	EndedAt    *time.Time         `json:"endedAt,omitempty"`
	Source     commonmodel.Source `json:"source"`
	Comments   string             `json:"comments"`
}

func NewServiceLineItemCreateEvent(aggregate eventstore.Aggregate, dataFields model.ServiceLineItemDataFields, source commonmodel.Source, createdAt, updatedAt, startedAt time.Time, endedAt *time.Time) (eventstore.Event, error) {
	eventData := ServiceLineItemCreateEvent{
		Tenant:     aggregate.GetTenant(),
		Billed:     dataFields.Billed.String(),
		Quantity:   dataFields.Quantity,
		Price:      dataFields.Price,
		Name:       dataFields.Name,
		ContractId: dataFields.ContractId,
		ParentId:   dataFields.ParentId,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		StartedAt:  startedAt,
		EndedAt:    endedAt,
		Source:     source,
		Comments:   dataFields.Comments,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ServiceLineItemCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ServiceLineItemCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ServiceLineItemCreateEvent")
	}

	return event, nil
}
