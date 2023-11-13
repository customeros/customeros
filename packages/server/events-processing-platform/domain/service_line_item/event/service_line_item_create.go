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
	Tenant      string             `json:"tenant" validate:"required"`
	Billed      string             `json:"billed"`
	Licenses    int32              `json:"licenses,omitempty"`
	Price       float32            `json:"price"`
	Description string             `json:"description"`
	ContractId  string             `json:"contractId"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	Source      commonmodel.Source `json:"source"`
}

func NewServiceLineItemCreateEvent(aggregate eventstore.Aggregate, dataFields model.ServiceLineItemDataFields, source commonmodel.Source, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ServiceLineItemCreateEvent{
		Tenant:      aggregate.GetTenant(),
		Billed:      dataFields.Billed.String(),
		Licenses:    dataFields.Licenses,
		Price:       dataFields.Price,
		Description: dataFields.Description,
		ContractId:  dataFields.ContractId,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Source:      source,
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
