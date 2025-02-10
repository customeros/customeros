package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var _ interfaces.AgentEvents = InvoiceStart{}

type InvoiceStart struct {
	AgentExecutionId string `json:"agentExecutionId"`
	ContractID       string `json:"contractId"`
	DryRun           bool   `json:"dryRun"`
	Preview          bool   `json:"preview"`
}
