package dto

type InvoicePaid struct {
	InvoiceID string `json:"invoiceId"`
	DryRun    bool   `json:"dryRun"`
}
