package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type OrganizationPlanCreateEvent struct {
	Tenant             string        `json:"tenant" validate:"required"`
	Name               string        `json:"name"`
	CreatedAt          time.Time     `json:"createdAt"`
	SourceFields       events.Source `json:"sourceFields"`
	OrganizationPlanId string        `json:"organizationPlanId"`
	MasterPlanId       string        `json:"masterPlanId"`
	OrganizationId     string        `json:"organizationId"`
}

func NewOrganizationPlanCreateEvent(aggregate eventstore.Aggregate, organizationPlanId, masterPlanId, organizationId, name string, sourceFields events.Source, createdAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationPlanCreateEvent{
		Tenant:             aggregate.GetTenant(),
		Name:               name,
		CreatedAt:          createdAt,
		SourceFields:       sourceFields,
		OrganizationPlanId: organizationPlanId,
		MasterPlanId:       masterPlanId,
		OrganizationId:     organizationId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationPlanCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationPlanCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationPlanCreateEvent")
	}

	return event, nil
}
