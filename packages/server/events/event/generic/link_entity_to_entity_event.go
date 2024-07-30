package generic

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event"
)

type LinkEntityWithEntity struct {
	event.BaseEvent
	WithEntityId     string           `json:"withEntityId"`
	WithEntityType   model.EntityType `json:"withEntityType"`
	RelationshipName string           `json:"relationshipName"`
	//TODO enhance with relationship properties
}

func (e LinkEntityWithEntity) GetBaseEvent() event.BaseEvent {
	return e.BaseEvent
}

func (e LinkEntityWithEntity) SetEntityId(entityId string) {
}
