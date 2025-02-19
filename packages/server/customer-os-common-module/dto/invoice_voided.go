package dto

import "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

type InvoiceVoided struct {
	InvoiceID string `json:"invoiceId"`
	DryRun    bool   `json:"dryRun"`
}

func (i InvoiceVoided) Name() enum.AgentListenerEvent {
	return enum.EventInvoiceVoided
}
