package notification

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
)

type notificationService struct {
	log                  logger.Logger
	postgresRepositories *postgres_repository.Repositories
	slackService         interfaces.SlackService
}

func NewNotificationService(log logger.Logger, postgresRepo *postgres_repository.Repositories, slackService interfaces.SlackService) interfaces.NotificationService {
	return &notificationService{
		log:                  log,
		postgresRepositories: postgresRepo,
		slackService:         slackService,
	}
}

func (s *notificationService) NotifySlackChannel(ctx context.Context, channelID string, message string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NotificationService.NotifySlackChannel")
	defer spans.Finish()

	err := s.slackService.SendMessageFromBot(ctx, channelID, message, true)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
