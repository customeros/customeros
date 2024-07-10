package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
)

type UpsertEmailToEntity struct {
	events.TenantBaseEvent

	EmailId  *string `json:"emailId"`
	RawEmail *string `json:"rawEmail"`

	EntityId   string            `json:"entityId"`
	EntityType events.EntityType `json:"entityType"`
}
