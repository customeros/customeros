package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/slack-go/slack"
	"time"
)

type SlackService interface {
	FetchUserIdsFromSlackChannel(ctx context.Context, channelId string, slackDtls SlackWorkspaceDtls) ([]string, error)
	FetchUserInfo(ctx context.Context, userId string, slackDtls SlackWorkspaceDtls) (*slack.User, error)
	FetchNewMessagesFromSlackChannel(ctx context.Context, channelId string, from, to time.Time, slackDtls SlackWorkspaceDtls) ([]slack.Message, error)
}

type slackService struct {
	cfg          *config.Config
	log          logger.Logger
	repositories *repository.Repositories
}

func NewSlackService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories) SlackService {
	return &slackService{
		cfg:          cfg,
		log:          log,
		repositories: repositories,
	}
}

func (s *slackService) FetchUserIdsFromSlackChannel(ctx context.Context, channelId string, slackDtls SlackWorkspaceDtls) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.FetchUserIdsFromSlackChannel")
	defer span.Finish()

	client := slack.New(slackDtls.token)

	users := make([]string, 0)

	var cursor = "" // initial empty cursor
	for {
		params := slack.GetUsersInConversationParameters{
			ChannelID: channelId,
			Cursor:    cursor,
		}
		members, cursor, err := client.GetUsersInConversation(&params)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		for _, member := range members {
			users = append(users, member)
		}

		if cursor == "" {
			break // no more pages
		}
	}

	return users, nil
}

func (s *slackService) FetchUserInfo(ctx context.Context, userId string, slackDtls SlackWorkspaceDtls) (*slack.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.FetchUserFromSlackChannel")
	defer span.Finish()

	client := slack.New(slackDtls.token)

	slackUser, err := client.GetUserInfo(userId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return slackUser, nil
}

func (s *slackService) FetchNewMessagesFromSlackChannel(ctx context.Context, channelId string, from, to time.Time, slackDtls SlackWorkspaceDtls) ([]slack.Message, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.FetchNewMessagesFromSlackChannel")
	defer span.Finish()
	span.LogFields(log.String("channelId", channelId), log.Object("from", from), log.Object("to", to))

	client := slack.New(slackDtls.token)

	messages := make([]slack.Message, 0)

	var cursor = "" // initial empty cursor
	for {
		params := slack.GetConversationHistoryParameters{
			ChannelID: channelId,
			Cursor:    cursor,
			Oldest:    toFloatTs(from),
			Latest:    toFloatTs(to),
			Inclusive: true,
		}
		response, err := client.GetConversationHistory(&params)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		for _, message := range response.Messages {
			messages = append(messages, message)
		}

		if response.ResponseMetaData.NextCursor == "" {
			break // no more pages
		}
		cursor = response.ResponseMetaData.NextCursor
	}

	return messages, nil
}

func toFloatTs(time time.Time) string {
	timestampMcs := float64(time.UnixNano()) / 1e9
	return fmt.Sprintf("%.3f", timestampMcs)
}
