package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type TenantBankAccountDeleteEvent struct {
	Tenant    string    `json:"tenant" validate:"required"`
	Id        string    `json:"id" validate:"required"`
	DeletedAt time.Time `json:"deletedAt"`
}

func NewTenantBankAccountDeleteEvent(aggregate eventstore.Aggregate, id string, deletedAt time.Time) (eventstore.Event, error) {
	eventData := TenantBankAccountDeleteEvent{
		Tenant:    aggregate.GetTenant(),
		Id:        id,
		DeletedAt: deletedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate TenantBankAccountDeleteEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, TenantDeleteBankAccountV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for TenantBankAccountDeleteEvent")
	}

	return event, nil
}
