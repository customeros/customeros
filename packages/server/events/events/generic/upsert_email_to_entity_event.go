package generic

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events"
)

type UpsertEmailToEntity struct {
	events.BaseEvent

	EmailId  *string `json:"emailId"`
	RawEmail *string `json:"rawEmail"`

	EntityId   string `json:"toEntityId"`
	EntityType string `json:"toEntityType"`
}
