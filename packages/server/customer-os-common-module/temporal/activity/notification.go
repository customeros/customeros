package activity

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/notifications"
)

func NotifyUserActivity(notification *notifications.NovuNotification, provider notifications.NotificationProvider) error {
	return provider.SendNotification(context.Background(), notification)
}
