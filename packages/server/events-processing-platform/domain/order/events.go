package order

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

const (
	OrderUpsertV1 = "V1_ORDER_UPSERT"
)

type OrderUpsertEvent struct {
	Tenant         string `json:"tenant" validate:"required"`
	OrganizationId string `json:"organizationId" validate:"required"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	ExternalSystem commonmodel.ExternalSystem `json:"externalSystem,omitempty"`
	SourceFields   events.Source              `json:"sourceFields"`

	ConfirmedAt *time.Time `json:"confirmedAt"`
	PaidAt      *time.Time `json:"paidAt"`
	FulfilledAt *time.Time `json:"fulfilledAt"`
	CanceledAt  *time.Time `json:"canceledAt"`
}

func NewOrderUpsertEvent(aggregate eventstore.Aggregate, sourceFields events.Source, externalSystem commonmodel.ExternalSystem, organizationId string, createdAt time.Time, updatedAt time.Time, confirmedAt, paidAt, fulfilledAt, canceledAt *time.Time) (eventstore.Event, error) {
	eventData := OrderUpsertEvent{
		Tenant:         aggregate.GetTenant(),
		SourceFields:   sourceFields,
		ExternalSystem: externalSystem,
		OrganizationId: organizationId,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		ConfirmedAt:    confirmedAt,
		PaidAt:         paidAt,
		FulfilledAt:    fulfilledAt,
		CanceledAt:     canceledAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrderUpsertEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrderUpsertV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrderUpsertEvent")
	}

	return event, nil
}
