package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type IntentDetected struct {
	EventName      enum.AgentListenerEvent
	Source         enum.Source
	SourceID       string            // Id of source record producing event
	IntentType     enum.IntentSignal `json:"intentType"`
	OrganizationID string            `json:"organizationId"`
	ContactID      string            `json:"contactId"`
	ContractID     string            `json:"contractId"`
	DryRun         bool              `json:"dryRun"`
	Preview        bool              `json:"preview"`
}
