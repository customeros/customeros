package notification

import (
	"context"
	"github.com/opentracing/opentracing-go"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
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

func (s *notificationService) NotifySlackChannel(ctx context.Context, tenant, channelID string, message *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NotificationService.NotifySlackChannel")
	tracing.SetDefaultServiceSpanTags(ctx, span)

	err := s.slackService.SendMessageFromBot(ctx, channelID, *message)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
