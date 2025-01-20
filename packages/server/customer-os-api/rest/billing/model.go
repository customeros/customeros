// @openapi 3.0.0
package billing

import (
	"time"
)

// InvoiceResponse represents a single invoice response
// @Description Response containing a single invoice's details
type InvoiceResponse struct {
	// The invoice information
	// required: false
	Invoice InvoiceRecord `json:"invoice,omitempty"`
}

// InvoicesResponse represents a collection of invoices
// @Description Response containing multiple invoices
type InvoicesResponse struct {
	// List of invoices
	// required: false
	Invoices []InvoiceRecord `json:"invoices,omitempty"`
}

// InvoiceRecord represents detailed invoice information
// @Description Detailed invoice information including payment details and status
type InvoiceRecord struct {
	// Unique identifier for the invoice
	// required: true
	// example: 123e4567-e89b-12d3-a456-426614174000
	// format: uuid
	ID string `json:"id"`

	// Invoice number or reference
	// required: true
	// example: INV-2024-001
	// minLength: 1
	Number string `json:"number"`

	// Date when the invoice payment is due
	// required: true
	// example: 2024-12-01T00:00:00Z
	// format: date-time
	DueDate time.Time `json:"dueDate"`

	// Current status of the invoice
	// required: true
	// example: PAID
	// enum: DRAFT,PENDING,PAID,OVERDUE,CANCELLED,VOID
	InvoiceStatus string `json:"invoiceStatus"`

	// Total amount due for the invoice
	// required: true
	// example: 1500.50
	// minimum: 0
	Amount float64 `json:"amount"`

	// Currency code for the invoice amount
	// required: true
	// example: USD
	// pattern: ^[A-Z]{3}$
	Currency string `json:"currency"`

	// URL where the invoice can be paid
	// required: false
	// example: https://payment.example.com/inv/12345
	// format: uri
	PaymentLink string `json:"paymentLink,omitempty"`

	// Public URL to access the invoice PDF
	// required: false
	// example: https://invoices.example.com/12345.pdf
	// format: uri
	PublicUrl string `json:"publicUrl,omitempty"`
}
