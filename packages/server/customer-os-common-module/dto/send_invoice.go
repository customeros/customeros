package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type SendInvoice struct {
	InvoiceId string `json:"invoiceId"`
}

func (e SendInvoice) Name() enum.AgentListenerEvent {
	return enum.EventSendInvoice
}
