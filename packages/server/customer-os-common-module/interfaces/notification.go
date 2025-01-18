package interfaces

import "context"

type NotificationService interface {
	NotifySlackChannel(ctx context.Context, tenant, channelID string, message *string) error
}
