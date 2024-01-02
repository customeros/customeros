package notifications

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

const (
	WorkflowIdTestFlow                      = "test-workflow"
	WorkflowIdOrgOwnerUpdateEmail           = "org-owner-update-email"
	WorkflowIdOrgOwnerUpdateAppNotification = "org-owner-update-in-app-notification"
)

type NotifiableUser struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	SubscriberID string `json:"subscriberId"` // must be unique uuid for user
}

type NotificationProvider interface {
	SendNotification(ctx context.Context, u *NotifiableUser, payload map[string]interface{}, workflowId string) error
	// SendInAppNotification(ctx context.Context, u *NotifiableUser, payload map[string]interface{}, eventId string) error
}

func NewNotificationProvider(log logger.Logger, apiKey string) NotificationProvider {
	return NewNovuProvider(log, apiKey)
}
