package events

import (
	"time"
)

type BaseEvent struct {
	CreatedAt time.Time `json:"createdAt"`
	AppSource string    `json:"appSource"`
	Source    string    `json:"source"`

	EventName string `json:"eventName" validate:"required"`

	Tenant     string     `json:"tenant" validate:"required"`
	EntityId   string     `json:"entityId"`
	EntityType EntityType `json:"entityType"`
}
