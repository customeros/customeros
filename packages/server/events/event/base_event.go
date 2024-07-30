package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"time"
)

type BaseEvent struct {
	CreatedAt      time.Time        `json:"createdAt"`
	AppSource      string           `json:"appSource"`
	Source         string           `json:"source"`
	LoggedInUserId string           `json:"loggedInUserId"`
	EventName      string           `json:"eventName" validate:"required"`
	Tenant         string           `json:"tenant" validate:"required"`
	EntityId       string           `json:"entityId" validate:"required"`
	EntityType     model.EntityType `json:"entityType" validate:"required"`
}

type BaseEventAccessor interface {
	GetBaseEvent() BaseEvent
	SetEntityId(entityId string)
}
