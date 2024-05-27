package notifications

const (
	WorkflowIdOrgOwnerUpdateEmailSubject           = "%s %s added you as an owner"
	WorkflowIdOrgOwnerUpdateAppNotificationSubject = "%s %s added you as an owner"
	WorkflowFailedWebhookSubject                   = "[Action Required] Webhook %s is offline"
	WorkflowInvoiceVoidedSubject                   = "Voided Invoice %s"
	WorkflowInvoicePaidSubject                     = "Paid Invoice %s from %s"
	WorkflowInvoicePaymentReceivedSubject          = "Payment Received for Invoice %s from %s"
	WorkflowInvoiceReadySubject                    = "New invoice %s"
	WorkflowReminderNotificationSubject            = "Reminder, %s"
)
