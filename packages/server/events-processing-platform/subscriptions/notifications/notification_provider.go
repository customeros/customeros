package notifications

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

const (
	EventIdTestFlow   = "test_flow"
	EventIdUserUpdate = "user_update"
)

type NotificationProvider interface {
	SendEmail(ctx context.Context, u *EmailableUser, payload map[string]interface{}, eventId string) error
	// TODO: SendInAppNotification(u *InAppNotifiableUser)
}

func NewNotificationProvider(log logger.Logger, apiKey string) NotificationProvider {
	return NewNovuProvider(log, apiKey)
}
