package dto

import "time"

type UpdateInvoice struct {
	Status                *string    `json:"status,omitempty"`
	PaymentLink           *string    `json:"paymentLink,omitempty"`
	PaymentLinkValidUntil *time.Time `json:"paymentLinkValidUntil,omitempty"`
}
