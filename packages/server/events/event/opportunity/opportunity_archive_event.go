package opportunity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/event"
)

type OpportunityArchiveEvent struct {
	event.BaseEvent
}

func (e OpportunityArchiveEvent) GetBaseEvent() event.BaseEvent {
	return e.BaseEvent
}

func (e *OpportunityArchiveEvent) SetEntityId(entityId string) {
}
