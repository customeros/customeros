package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"time"
)

const pageSize = 200

type SlackService interface {
	FetchChannelInfo(ctx context.Context, token, channelId string) (*slack.Channel, error)
	AuthTest(ctx context.Context, token string) (*slack.AuthTestResponse, error)
	FetchUserIdsFromSlackChannel(ctx context.Context, token, channelId string) ([]string, error)
	FetchUserInfo(ctx context.Context, token, userId string) (*slack.User, error)
	FetchNewMessagesFromSlackChannel(ctx context.Context, token, channelId string, from, to time.Time) ([]slack.Message, error)
	FetchMessagesFromSlackChannelWithReplies(ctx context.Context, token, channelId string, to time.Time, lookbackWindow int) ([]slack.Message, error)
	FetchNewThreadMessages(ctx context.Context, token, channelId, parentTs string, from, to time.Time) ([]slack.Message, error)
	GetMessagePermalink(ctx context.Context, token, channelId, messageTs string) (string, error)
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

func (s *slackService) FetchChannelInfo(ctx context.Context, token, channelId string) (*slack.Channel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.FetchChannelInfo")
	defer span.Finish()

	client := slack.New(token)

	channel, err := client.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID:         channelId,
		IncludeLocale:     false,
		IncludeNumMembers: false,
	})
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	return channel, nil
}

func (s *slackService) AuthTest(ctx context.Context, token string) (*slack.AuthTestResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.AuthTest")
	defer span.Finish()

	client := slack.New(token)

	authTest, err := client.AuthTest()
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	return authTest, nil
}

func (s *slackService) FetchUserIdsFromSlackChannel(ctx context.Context, token, channelId string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.FetchUserIdsFromSlackChannel")
	defer span.Finish()

	client := slack.New(token)

	users := make([]string, 0)

	var cursor = "" // initial empty cursor
	for {
		params := slack.GetUsersInConversationParameters{
			ChannelID: channelId,
			Cursor:    cursor,
			Limit:     pageSize,
		}
		members, cursor, err := client.GetUsersInConversation(&params)
		if err != nil {
			if rlErr, ok := err.(*slack.RateLimitedError); ok {
				wait := rlErr.RetryAfter
				select {
				case <-time.After(wait):
					// retry after delay
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			} else {
				tracing.TraceErr(span, err)
				return nil, err
			}
		} else {
			for _, member := range members {
				users = append(users, member)
			}
			if cursor == "" {
				break // no more pages
			}
			select {
			case <-ctx.Done():
				break
			default:
			}
		}
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return users, nil
}

func (s *slackService) FetchUserInfo(ctx context.Context, token, userId string) (*slack.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.FetchUserInfo")
	defer span.Finish()
	span.LogFields(log.String("slackUserId", userId))

	client := slack.New(token)

	slackUser, err := client.GetUserInfo(userId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return slackUser, nil
}

func (s *slackService) FetchNewMessagesFromSlackChannel(ctx context.Context, token, channelId string, from, to time.Time) ([]slack.Message, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.FetchNewMessagesFromSlackChannel")
	defer span.Finish()
	span.LogFields(log.String("channelId", channelId), log.Object("from", from), log.Object("to", to))

	client := slack.New(token)

	messages := make([]slack.Message, 0, pageSize)

	var cursor = "" // initial empty cursor
	for {
		params := slack.GetConversationHistoryParameters{
			ChannelID: channelId,
			Cursor:    cursor,
			Oldest:    toFloatTs(from),
			Latest:    toFloatTs(to),
			Inclusive: true,
			Limit:     pageSize,
		}
		page, err := client.GetConversationHistory(&params)
		if err != nil {
			if rlErr, ok := err.(*slack.RateLimitedError); ok {
				wait := rlErr.RetryAfter
				select {
				case <-time.After(wait):
					// retry after delay
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			} else {
				tracing.TraceErr(span, err)
				return nil, err
			}
		} else {
			messages = append(messages, page.Messages...)
			cursor = page.ResponseMetaData.NextCursor
			if page.HasMore == false || cursor == "" {
				break // no more pages
			}
			select {
			case <-ctx.Done():
				break
			default:
			}
		}
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return messages, nil
}

func (s *slackService) FetchMessagesFromSlackChannelWithReplies(ctx context.Context, token, channelId string, to time.Time, lookbackWindow int) ([]slack.Message, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.FetchMessagesFromSlackChannelWithReplies")
	defer span.Finish()
	span.LogFields(log.String("channelId", channelId), log.Int("lookBackWindowDays", lookbackWindow), log.Object("to", to))

	client := slack.New(token)
	messages := make([]slack.Message, 0, pageSize)
	from := utils.Now().AddDate(0, 0, 0-lookbackWindow)

	var cursor = "" // initial empty cursor
	for {
		params := slack.GetConversationHistoryParameters{
			ChannelID: channelId,
			Cursor:    cursor,
			Oldest:    toFloatTs(from),
			Latest:    toFloatTs(to),
			Inclusive: true,
			Limit:     pageSize,
		}
		page, err := client.GetConversationHistory(&params)
		if err != nil {
			if rlErr, ok := err.(*slack.RateLimitedError); ok {
				wait := rlErr.RetryAfter
				select {
				case <-time.After(wait):
					// retry after delay
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			} else {
				tracing.TraceErr(span, err)
				return nil, err
			}
		} else {
			for _, msg := range page.Messages {
				// return only messages with replies
				if msg.ThreadTimestamp != "" {
					messages = append(messages, msg)
				}
			}
			cursor = page.ResponseMetaData.NextCursor
			if page.HasMore == false || cursor == "" {
				break // no more pages
			}
			select {
			case <-ctx.Done():
				break
			default:
			}
		}
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return messages, nil
}

func (s *slackService) FetchNewThreadMessages(ctx context.Context, token, channelId, parentTs string, from, to time.Time) ([]slack.Message, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.FetchNewThreadMessages")
	defer span.Finish()
	span.LogFields(log.String("channelId", channelId), log.String("parentMessageTs", parentTs), log.Object("from", from), log.Object("to", to))

	client := slack.New(token)

	var messages []slack.Message
	messages = make([]slack.Message, 0, pageSize)

	cursor := ""
	for {
		params := slack.GetConversationRepliesParameters{
			ChannelID: channelId,
			Timestamp: parentTs,
			Cursor:    cursor,
			Oldest:    toFloatTs(from),
			Latest:    toFloatTs(to),
			Inclusive: true,
			Limit:     pageSize,
		}
		pageMsgs, hasMore, cursor, err := client.GetConversationRepliesContext(ctx, &params)
		if err != nil {
			if rlErr, ok := err.(*slack.RateLimitedError); ok {
				wait := rlErr.RetryAfter
				select {
				case <-time.After(wait):
					// retry after delay
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			} else {
				tracing.TraceErr(span, err)
				return nil, err
			}
		} else {
			for _, message := range pageMsgs {
				if message.ThreadTimestamp != message.Timestamp {
					messages = append(messages, message)
				}
			}
			if hasMore == false || cursor == "" {
				break // no more pages
			}
			select {
			case <-ctx.Done():
				break
			default:
			}
		}
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return messages, nil
}

func (s *slackService) GetMessagePermalink(ctx context.Context, token, channelId, messageTs string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.GetMessagePermalink")
	defer span.Finish()
	span.LogKV("channelId", channelId, "messageTimestamp", messageTs)

	client := slack.New(token)

	params := slack.PermalinkParameters{
		Channel: channelId,
		Ts:      messageTs,
	}

	retryCount, maxRetryCount := 0, 5
	for {
		permalink, err := client.GetPermalinkContext(ctx, &params)
		if err != nil {
			if retryCount > maxRetryCount {
				return "", errors.New("max retry count reached")
			}
			retryCount++
			if rlErr, ok := err.(*slack.RateLimitedError); ok {
				wait := rlErr.RetryAfter
				select {
				case <-time.After(wait):
					// retry after delay
				case <-ctx.Done():
					return "", ctx.Err()
				}
			} else {
				tracing.TraceErr(span, err)
				return "", err
			}
		} else {
			return permalink, nil
		}
	}
}

func toFloatTs(time time.Time) string {
	timestampMcs := float64(time.UnixNano()) / 1e9
	return fmt.Sprintf("%.3f", timestampMcs)
}
