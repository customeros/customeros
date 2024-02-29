package notifications

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/aws_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
)

const (
	WorkflowIdTestFlow                      = "test-workflow"
	WorkflowIdOrgOwnerUpdateEmail           = "org-owner-update-email"
	WorkflowIdOrgOwnerUpdateAppNotification = "org-owner-update-in-app-notification"
	WorkflowInvoicePaid                     = "invoice-paid"
	WorkflowInvoiceReady                    = "invoice-ready"
	WorkflowInvoiceVoided                   = "invoice-voided"
	WorkflowFailedWebhook                   = "failed-webhook"
)

var REQUIRED_TEMPLATE_VALUES = map[string][]string{
	WorkflowIdOrgOwnerUpdateEmail: {
		"{{userFirstName}}",
		"{{actorFirstName}}",
		"{{actorLastName}}",
		"{{orgName}}",
		"{{orgLink}}",
	},
	WorkflowFailedWebhook: {
		"{{userFirstName}}",
		"{{tenant}}",
		"{{webhookUrl}}",
	},
}

type NotifiableUser struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	SubscriberID string `json:"subscriberId"` // must be unique uuid for user
}

type NovuNotification struct {
	WorkflowId   string
	TemplateData map[string]string

	To      *NotifiableUser
	Subject string
	Payload map[string]interface{}
}

type NotificationProvider interface {
	SendNotification(ctx context.Context, notification *NovuNotification) error
}

func NewNovuNotificationProvider(log logger.Logger, apiKey string, s3Client aws_client.S3ClientI) NotificationProvider {
	return newNovuProvider(log, apiKey, s3Client)
}
