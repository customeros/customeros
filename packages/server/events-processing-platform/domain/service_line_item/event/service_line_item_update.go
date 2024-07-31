package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type ServiceLineItemUpdateEvent struct {
	Tenant    string        `json:"tenant" validate:"required"`
	Name      string        `json:"name"`
	Quantity  int64         `json:"quantity,omitempty" validate:"min=0"`
	Price     float64       `json:"price,omitempty"`
	UpdatedAt time.Time     `json:"updatedAt"`
	Billed    string        `json:"billed"`
	Source    common.Source `json:"source"`
	Comments  string        `json:"comments"`
	VatRate   float64       `json:"vatRate"`
	StartedAt *time.Time    `json:"startedAt,omitempty"`
}

func NewServiceLineItemUpdateEvent(aggregate eventstore.Aggregate, dataFields model.ServiceLineItemDataFields, source common.Source, updatedAt time.Time, startedAt *time.Time) (eventstore.Event, error) {
	eventData := ServiceLineItemUpdateEvent{
		Tenant:    aggregate.GetTenant(),
		Name:      dataFields.Name,
		Quantity:  dataFields.Quantity,
		Price:     dataFields.Price,
		UpdatedAt: updatedAt,
		Billed:    dataFields.Billed.String(),
		Source:    source,
		Comments:  dataFields.Comments,
		VatRate:   dataFields.VatRate,
		StartedAt: utils.ToDatePtr(startedAt),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ServiceLineItemUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ServiceLineItemUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ServiceLineItemUpdateEvent")
	}
	return event, nil
}
