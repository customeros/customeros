package webhook

type WebhookEvent string

const (
	// WebhookEventInvoiceFinalized is the event name for invoice finalized
	WebhookEventInvoiceFinalized WebhookEvent = "invoice.finalized"
)

func (e WebhookEvent) String() string {
	return string(e)
}
