package interfaces

import "context"

type NotificationService interface {
	NotifySlackChannel(ctx context.Context, channelID string, message string) error
}
