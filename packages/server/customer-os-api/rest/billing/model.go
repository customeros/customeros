package billing

import "time"

// Invoice represents the structure of an invoice
// @Description Invoice details
type InvoiceResponse struct {
	// Status indicates the result of the action
	// Example: success
	Status string `json:"status,omitempty" example:"success"`

	// Message provides additional information about the action
	// Example: Invoices retrieved successfully
	Message string `json:"message,omitempty" example:"Invoices retrieved successfully"`

	// ID is the unique identifier for the invoice, uuid format
	// Example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`

	// Number represents the invoice number
	// Example: ABC-12345
	Number string `json:"number" example:"ABC-12345"`

	// DueDate represents the date the invoice is due
	// Example: 2024-12-01T00:00:00Z
	DueDate time.Time `json:"dueDate" example:"2024-12-01T00:00:00Z"`

	// Status represents the payment status of the invoice
	// Example: PAID
	InvoiceStatus string `json:"invoiceStatus" example:"PAID"`

	// Amount represents the total amount due for the invoice
	// Example: 1500.50
	Amount float64 `json:"amount" example:"1500.50"`

	// Currency represents the currency used for the invoice
	// Example: USD
	Currency string `json:"currency" example:"USD"`

	// PaymentLink represents the URL where the invoice can be paid
	// Example: https://example.com/payments/12345
	PaymentLink string `json:"paymentLink" example:"https://example.com/payments/12345"`

	// PublicUrl represents the public URL where the PDF version of the invoice can be accessed
	// Example: https://example.com/invoices/12345.pdf
	PublicUrl string `json:"publicUrl" example:"https://example.com/invoices/12345.pdf"`
}

// InvoicesResponse defines the response structure for multiple invoices in the response
// @Description Response body for all invoices details
type InvoicesResponse struct {
	// Status indicates the result of the action
	// Example: success
	Status string `json:"status,omitempty" example:"success"`

	// Message provides additional information about the action
	// Example: Invoices retrieved successfully
	Message string `json:"message,omitempty" example:"Invoices retrieved successfully"`

	Invoices []InvoiceResponse `json:"invoices"`
}
