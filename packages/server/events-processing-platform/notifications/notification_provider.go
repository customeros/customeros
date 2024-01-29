package notifications

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

const (
	WorkflowIdTestFlow                      = "test-workflow"
	WorkflowIdOrgOwnerUpdateEmail           = "org-owner-update-email"
	WorkflowIdOrgOwnerUpdateAppNotification = "org-owner-update-in-app-notification"
	WorkflowInvoicePaid                     = "invoice-paid"
	WorkflowInvoiceReady                    = "invoice-ready"
)

type NotifiableUser struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	SubscriberID string `json:"subscriberId"` // must be unique uuid for user
}

type NotificationProvider interface {
	SendNotification(ctx context.Context, u *NotifiableUser, payload map[string]interface{}, workflowId string) error
	LoadEmailBody(ctx context.Context, workflowId string) error
	Template(ctx context.Context, replace map[string]string) (string, error)
	GetRawContent() string
}

func NewNovuNotificationProvider(log logger.Logger, apiKey, templatePath string) NotificationProvider {
	return newNovuProvider(log, apiKey, templatePath)
}

func NewPostmarkNotificationProvider(log logger.Logger, serverToken, accountToken string) NotificationProvider {
	return newPostmarkProvider(log, serverToken, accountToken)
}
