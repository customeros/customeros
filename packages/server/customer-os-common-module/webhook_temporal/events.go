package webhook

type WebhookEvent string

const (
	// WebhookEventInvoiceFinalized is the event name for invoice finalized
	WebhookEventInvoiceFinalized  WebhookEvent = "invoice.finalized"
	WebhookEventInvoiceStatusPaid WebhookEvent = "invoice.status.paid"
)

func (e WebhookEvent) String() string {
	return string(e)
}
