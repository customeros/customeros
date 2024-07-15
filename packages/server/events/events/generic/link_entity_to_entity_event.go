package generic

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
)

type LinkEntityWithEntity struct {
	events.BaseEvent

	WithEntityId   string           `json:"withEntityId"`
	WithEntityType model.EntityType `json:"withEntityType"`

	RelationshipName string `json:"relationshipName"`

	//todo enhance with relationship properties
}

func (e LinkEntityWithEntity) GetBaseEvent() events.BaseEvent {
	return e.BaseEvent
}

func (e *LinkEntityWithEntity) SetEntityId(entityId string) {
	e.EntityId = entityId
}
