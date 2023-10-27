package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	InteractionEventAggregateType eventstore.AggregateType = "interaction_event"
)

type InteractionEventAggregate struct {
	*aggregate.CommonTenantIdAggregate
	InteractionEvent *model.InteractionEvent
}

func NewInteractionEventAggregateWithTenantAndID(tenant, id string) *InteractionEventAggregate {
	interactionEventAggregate := InteractionEventAggregate{}
	interactionEventAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(InteractionEventAggregateType, tenant, id)
	interactionEventAggregate.SetWhen(interactionEventAggregate.When)
	interactionEventAggregate.InteractionEvent = &model.InteractionEvent{}
	interactionEventAggregate.Tenant = tenant
	return &interactionEventAggregate
}

func (a *InteractionEventAggregate) When(evt eventstore.Event) error {

	switch evt.GetEventType() {
	case event.InteractionEventCreateV1:
		return a.onInteractionEventCreate(evt)
	case event.InteractionEventUpdateV1:
		return a.onInteractionEventUpdate(evt)
	case event.InteractionEventRequestSummaryV1,
		event.InteractionEventRequestActionItemsV1:
		return nil
	case event.InteractionEventReplaceSummaryV1:
		return a.onSummaryReplace(evt)
	case event.InteractionEventReplaceActionItemsV1:
		return a.onActionItemsReplace(evt)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *InteractionEventAggregate) onSummaryReplace(evt eventstore.Event) error {
	var eventData event.InteractionEventReplaceSummaryEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.InteractionEvent.Summary = eventData.Summary
	a.InteractionEvent.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *InteractionEventAggregate) onActionItemsReplace(evt eventstore.Event) error {
	var eventData event.InteractionEventReplaceActionItemsEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.InteractionEvent.UpdatedAt = eventData.UpdatedAt
	if len(eventData.ActionItems) > 0 {
		a.InteractionEvent.ActionItems = eventData.ActionItems
	}
	return nil
}

func (a *InteractionEventAggregate) onInteractionEventCreate(evt eventstore.Event) error {
	var eventData event.InteractionEventCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.InteractionEvent.ID = a.ID
	a.InteractionEvent.Tenant = a.Tenant
	a.InteractionEvent.Content = eventData.Content
	a.InteractionEvent.ContentType = eventData.ContentType
	a.InteractionEvent.Channel = eventData.Channel
	a.InteractionEvent.ChannelData = eventData.ChannelData
	a.InteractionEvent.EventType = eventData.EventType
	a.InteractionEvent.Identifier = eventData.Identifier
	a.InteractionEvent.BelongsToSessionId = eventData.BelongsToSessionId
	a.InteractionEvent.BelongsToIssueId = eventData.BelongsToIssueId
	a.InteractionEvent.Source = cmnmod.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.Source,
		AppSource:     eventData.AppSource,
	}
	a.InteractionEvent.CreatedAt = eventData.CreatedAt
	a.InteractionEvent.UpdatedAt = eventData.UpdatedAt
	if eventData.ExternalSystem.Available() {
		a.InteractionEvent.ExternalSystems = []cmnmod.ExternalSystem{eventData.ExternalSystem}
	}
	a.InteractionEvent.Hide = eventData.Hide
	return nil
}

func (a *InteractionEventAggregate) onInteractionEventUpdate(evt eventstore.Event) error {
	var eventData event.InteractionEventUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if eventData.Source == constants.SourceOpenline {
		a.InteractionEvent.Source.SourceOfTruth = eventData.Source
	}
	if eventData.Source != a.InteractionEvent.Source.SourceOfTruth && a.InteractionEvent.Source.SourceOfTruth == constants.SourceOpenline {
		if a.InteractionEvent.Content == "" {
			a.InteractionEvent.Content = eventData.Content
		}
		if a.InteractionEvent.ContentType == "" {
			a.InteractionEvent.ContentType = eventData.ContentType
		}
		if a.InteractionEvent.Channel == "" {
			a.InteractionEvent.Channel = eventData.Channel
		}
		if a.InteractionEvent.ChannelData == "" {
			a.InteractionEvent.ChannelData = eventData.ChannelData
		}
		if a.InteractionEvent.Identifier == "" {
			a.InteractionEvent.Identifier = eventData.Identifier
		}
		if a.InteractionEvent.EventType == "" {
			a.InteractionEvent.EventType = eventData.EventType
		}
	} else {
		a.InteractionEvent.Content = eventData.Content
		a.InteractionEvent.ContentType = eventData.ContentType
		a.InteractionEvent.Channel = eventData.Channel
		a.InteractionEvent.ChannelData = eventData.ChannelData
		a.InteractionEvent.Identifier = eventData.Identifier
		a.InteractionEvent.EventType = eventData.EventType
		a.InteractionEvent.Hide = eventData.Hide
	}
	a.InteractionEvent.UpdatedAt = eventData.UpdatedAt
	return nil
}
