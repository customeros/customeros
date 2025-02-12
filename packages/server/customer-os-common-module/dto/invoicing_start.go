package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type InvoiceContract struct {
	ContractId string `json:"contractId"`
}

func (e InvoiceContract) Name() enum.AgentListenerEvent {
	return enum.EventStartInvoiceRun
}
