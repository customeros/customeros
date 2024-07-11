package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type OrganizationRequestEnrich struct {
	Tenant      string    `json:"tenant" validate:"required"`
	Website     string    `json:"website"`
	RequestedAt time.Time `json:"requestedAt"`
}

func NewOrganizationRequestEnrich(aggregate eventstore.Aggregate, website string) (eventstore.Event, error) {
	eventData := OrganizationRequestEnrich{
		Tenant:      aggregate.GetTenant(),
		Website:     website,
		RequestedAt: utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationRequestEnrich")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationRequestEnrichV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationRequestEnrich")
	}
	return event, nil
}
