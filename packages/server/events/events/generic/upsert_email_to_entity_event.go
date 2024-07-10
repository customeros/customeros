package generic

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
)

const (
	UpsertEmailToEntityV1 = "V1_UPSERT_EMAIL_TO_ENTITY"
)

type UpsertEmailToEntity struct {
	events.BaseEvent

	EmailId  *string `json:"emailId"`
	RawEmail *string `json:"rawEmail"`
}
