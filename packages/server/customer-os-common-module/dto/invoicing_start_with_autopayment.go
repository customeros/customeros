package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type InvoiceContractWithAutopayment struct {
	ContractId string `json:"contractId"`
}

func (e InvoiceContractWithAutopayment) Name() enum.AgentListenerEvent {
	return enum.EventStartInvoiceRunWithAutopayment
}
