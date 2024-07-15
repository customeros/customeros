package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type OrganizationPlanMilestoneCreateEvent struct {
	Tenant             string        `json:"tenant" validate:"required"`
	MilestoneId        string        `json:"milestoneId" validate:"required"`
	Name               string        `json:"name"`
	Order              int64         `json:"order" validate:"gte=0"`
	DueDate            time.Time     `json:"dueDate" validate:"gte=0"`
	CreatedAt          time.Time     `json:"createdAt"`
	Items              []string      `json:"items"`
	SourceFields       common.Source `json:"sourceFields"`
	Optional           bool          `json:"optional"`
	OrganizationPlanId string        `json:"organizationPlanId" validate:"required"`
	Adhoc              bool          `json:"adhoc"`
}

func NewOrganizationPlanMilestoneCreateEvent(aggregate eventstore.Aggregate, organizationPlanId, milestoneId, name string, order int64, items []string, optional, adhoc bool, sourceFields common.Source, createdAt, dueDate time.Time) (eventstore.Event, error) {
	eventData := OrganizationPlanMilestoneCreateEvent{
		Tenant:             aggregate.GetTenant(),
		MilestoneId:        milestoneId,
		Name:               name,
		CreatedAt:          createdAt,
		Order:              order,
		Items:              items,
		SourceFields:       sourceFields,
		Optional:           optional,
		OrganizationPlanId: organizationPlanId,
		DueDate:            dueDate,
		Adhoc:              adhoc,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationPlanMilestoneCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationPlanMilestoneCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationPlanMilestoneCreateEvent")
	}

	return event, nil
}
