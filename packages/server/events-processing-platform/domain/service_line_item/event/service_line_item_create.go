package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type ServiceLineItemCreateEvent struct {
	Tenant            string        `json:"tenant" validate:"required"`
	Billed            string        `json:"billed"`
	Quantity          int64         `json:"quantity,omitempty" validate:"min=0"`
	Price             float64       `json:"price,omitempty" validate:"min=0"`
	Name              string        `json:"name"`
	ContractId        string        `json:"contractId" validate:"required"`
	ParentId          string        `json:"parentId" validate:"required"`
	PreviousVersionId string        `json:"previousVersionId"`
	CreatedAt         time.Time     `json:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt"`
	StartedAt         time.Time     `json:"startedAt"`
	EndedAt           *time.Time    `json:"endedAt,omitempty"`
	Source            common.Source `json:"source"`
	Comments          string        `json:"comments,omitempty"`
	VatRate           float64       `json:"vatRate"`
}

func NewServiceLineItemCreateEvent(aggregate eventstore.Aggregate, dataFields model.ServiceLineItemDataFields, source common.Source, createdAt, updatedAt, startedAt time.Time, endedAt *time.Time, previousVersionId string) (eventstore.Event, error) {
	eventData := ServiceLineItemCreateEvent{
		Tenant:            aggregate.GetTenant(),
		Billed:            dataFields.Billed.String(),
		Quantity:          dataFields.Quantity,
		Price:             dataFields.Price,
		Name:              dataFields.Name,
		ContractId:        dataFields.ContractId,
		ParentId:          dataFields.ParentId,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		StartedAt:         utils.ToDate(startedAt),
		EndedAt:           endedAt,
		Source:            source,
		Comments:          dataFields.Comments,
		VatRate:           dataFields.VatRate,
		PreviousVersionId: previousVersionId,
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
