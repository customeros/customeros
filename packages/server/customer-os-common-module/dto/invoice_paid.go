package dto

import "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

type InvoicePaid struct {
	InvoiceID string `json:"invoiceId"`
	DryRun    bool   `json:"dryRun"`
}

func (i InvoicePaid) Name() enum.AgentListenerEvent {
	return enum.EventInvoicePaid
}
