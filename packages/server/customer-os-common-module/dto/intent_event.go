package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type IntentEvent struct {
	EventName      enum.AgentListenerEvent `json:"eventName"`
	Source         enum.Source             `json:"source"`
	SourceID       string                  `json:"sourceId"`
	Tenant         string
	IntentType     enum.IntentSignal `json:"intentType"`
	OrganizationID string            `json:"organizationId"`
	ContactID      string            `json:"contactId"`
}
