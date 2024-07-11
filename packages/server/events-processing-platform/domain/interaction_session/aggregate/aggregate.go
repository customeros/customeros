package aggregate

import (
	events2 "github.com/openline-ai/openline-customer-os/packages/server/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"strings"
)

const (
	InteractionSessionAggregateType eventstore.AggregateType = "interaction_session"
)

type InteractionSessionAggregate struct {
	*eventstore.CommonTenantIdAggregate
	InteractionSession *model.InteractionSession
}

func NewInteractionSessionAggregateWithTenantAndID(tenant, id string) *InteractionSessionAggregate {
	interactionEventAggregate := InteractionSessionAggregate{}
	interactionEventAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(InteractionSessionAggregateType, tenant, id)
	interactionEventAggregate.SetWhen(interactionEventAggregate.When)
	interactionEventAggregate.InteractionSession = &model.InteractionSession{}
	interactionEventAggregate.Tenant = tenant
	return &interactionEventAggregate
}

func (a *InteractionSessionAggregate) When(evt eventstore.Event) error {

	switch evt.GetEventType() {
	case event.InteractionSessionCreateV1:
		return a.onInteractionSessionCreate(evt)
	default:
		if strings.HasPrefix(evt.GetEventType(), events2.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *InteractionSessionAggregate) onInteractionSessionCreate(evt eventstore.Event) error {
	var eventData event.InteractionSessionCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.InteractionSession.ID = a.ID
	a.InteractionSession.Tenant = a.Tenant
	a.InteractionSession.Channel = eventData.Channel
	a.InteractionSession.ChannelData = eventData.ChannelData
	a.InteractionSession.Status = eventData.Status
	a.InteractionSession.Type = eventData.Type
	a.InteractionSession.Name = eventData.Name
	a.InteractionSession.Identifier = eventData.Identifier
	a.InteractionSession.Source = events.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.Source,
		AppSource:     eventData.AppSource,
	}
	a.InteractionSession.CreatedAt = eventData.CreatedAt
	a.InteractionSession.UpdatedAt = eventData.UpdatedAt
	if eventData.ExternalSystem.Available() {
		a.InteractionSession.ExternalSystems = []cmnmod.ExternalSystem{eventData.ExternalSystem}
	}
	return nil
}
