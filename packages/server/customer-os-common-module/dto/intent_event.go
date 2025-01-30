package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type IntentDetected struct {
	EventName      enum.AgentListenerEvent
	Source         enum.Source
	SourceID       string // Id of source record producing event
	Tenant         string
	IntentType     enum.IntentSignal
	OrganizationID string
	ContactID      string
}
