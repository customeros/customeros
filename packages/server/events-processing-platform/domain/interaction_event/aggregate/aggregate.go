package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	InteractionEventAggregateType eventstore.AggregateType = "interaction_event"
)

type InteractionEventAggregate struct {
	*aggregate.CommonTenantIdAggregate
	InteractionEvent *models.InteractionEvent
}

func NewInteractionEventAggregateWithTenantAndID(tenant, id string) *InteractionEventAggregate {
	interactionEventAggregate := InteractionEventAggregate{}
	interactionEventAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(InteractionEventAggregateType, tenant, id)
	interactionEventAggregate.SetWhen(interactionEventAggregate.When)
	interactionEventAggregate.InteractionEvent = &models.InteractionEvent{}
	return &interactionEventAggregate
}

func (a *InteractionEventAggregate) When(event eventstore.Event) error {

	switch event.GetEventType() {
	case events.InteractionEventRequestSummaryV1,
		events.InteractionEventRequestActionItemsV1:
		return nil
	case events.InteractionEventReplaceSummaryV1:
		return a.onSummaryReplace(event)
	case events.InteractionEventReplaceActionItemsV1:
		return a.onActionItemsReplace(event)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		return err
	}
}

func (a *InteractionEventAggregate) onSummaryReplace(event eventstore.Event) error {
	var eventData events.InteractionEventReplaceSummaryEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.InteractionEvent.Summary = eventData.Summary
	a.InteractionEvent.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *InteractionEventAggregate) onActionItemsReplace(event eventstore.Event) error {
	var eventData events.InteractionEventReplaceActionItemsEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.InteractionEvent.UpdatedAt = eventData.UpdatedAt
	if len(eventData.ActionItems) > 0 {
		a.InteractionEvent.ActionItems = eventData.ActionItems
	}
	return nil
}
