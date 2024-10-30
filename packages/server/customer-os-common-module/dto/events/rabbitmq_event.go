package events

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"

type Event struct {
	Event    EventDetails  `json:"event"`
	Metadata EventMetadata `json:"metadata"`
}

type EventDetails struct {
	Id         string           `json:"id"`
	Tenant     string           `json:"tenant"`
	EntityId   string           `json:"entityId"`
	EntityType model.EntityType `json:"entityType"`
	EventType  string           `json:"eventType"`
	Data       interface{}      `json:"data"`
}

type EventMetadata struct {
	UberTraceId string `json:"uber-trace-id"`
	AppSource   string `json:"appSource"`
	UserId      string `json:"userId"`
	UserEmail   string `json:"userEmail"`
	Timestamp   string `json:"timestamp"`
}
