package data_fields

import "time"

type InvoiceFields struct {
	DryRun           bool       `json:"dryRun"`
	Preview          bool       `json:"preview"`
	InvoiceStartDate *time.Time `json:"invoiceStartDate"`
	InvoiceEndDate   *time.Time `json:"invoiceEndDate"`
	InvoiceNumber    string     `json:"invoiceNumber"`
}
