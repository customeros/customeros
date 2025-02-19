package dto

import "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

type PastDueInvoice struct {
	InvoiceId string `json:"invoiceId"`
}

func (i PastDueInvoice) Name() enum.AgentListenerEvent {
	return enum.EventInvoicePastDue
}
