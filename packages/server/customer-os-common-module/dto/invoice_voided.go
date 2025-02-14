package dto

type InvoiceVoided struct {
	InvoiceID string `json:"invoiceId"`
	DryRun    bool   `json:"dryRun"`
}
