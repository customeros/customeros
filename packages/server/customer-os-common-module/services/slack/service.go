package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type slackService struct {
	log      logger.Logger
	cfg      *config.GlobalConfig
	postgres *repository.Repositories
}

func NewSlackService(log logger.Logger, config *config.GlobalConfig, postgres *repository.Repositories) interfaces.SlackService {
	return &slackService{
		log:      log,
		cfg:      config,
		postgres: postgres,
	}
}

func (s *slackService) GetSlackChannels(ctx context.Context, tenant string) ([]*postgresEntity.SlackChannel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.GetSlackChannels")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, "service")

	nodes, err := s.postgres.SlackChannelRepository.GetSlackChannels(ctx, tenant)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (s *slackService) GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) ([]*postgresEntity.SlackChannel, int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.GetSlackChannels")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, "service")

	channels, totalCount, err := s.postgres.SlackChannelRepository.GetPaginatedSlackChannels(ctx, tenant, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return channels, totalCount, nil
}

func (s *slackService) StoreSlackChannel(ctx context.Context, tenant, source, channelId, channelName string, organizationId *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.StoreSlackChannel")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, "service")
	span.LogFields(log.String("channelId", channelId))
	span.LogFields(log.String("channelName", channelName))
	span.LogFields(log.String("organizationId", utils.IfNotNilString(organizationId)))

	existing, err := s.postgres.SlackChannelRepository.GetSlackChannel(ctx, tenant, channelId)
	if err != nil {
		return err
	}

	if existing == nil {
		now := time.Now()
		slackChannel := postgresEntity.SlackChannel{
			CreatedAt:      now,
			UpdatedAt:      now,
			TenantName:     tenant,
			ChannelId:      channelId,
			ChannelName:    channelName,
			OrganizationId: organizationId,
			Source:         source,
		}
		return s.postgres.SlackChannelRepository.CreateSlackChannel(ctx, &slackChannel)
	}
	if existing != nil {
		if organizationId != nil {
			return s.postgres.SlackChannelRepository.UpdateSlackChannelOrganization(ctx, existing.ID, *organizationId)
		} else if channelName != "" {
			return s.postgres.SlackChannelRepository.UpdateSlackChannelName(ctx, existing.ID, channelName)
		}
	}

	return nil
}

func (s *slackService) Notify(ctx context.Context, tenant, channelID string, message *string) error {
	span, ctx := tracing.StartTracerSpan(ctx, "Slack.Notify")
	defer span.Finish()

	err := s.sendSlackMessage(ctx, tenant, channelID, *message)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *slackService) sendSlackMessage(ctx context.Context, tenant, channel, blocks string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingService.sendSlackMessage")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("channel", channel))

	// Create HTTP client
	client := &http.Client{}

	requestBody := map[string]interface{}{
		"channel":      channel,
		"unfurl_links": false,
		"unfurl_media": false,
		"blocks":       blocks,
	}

	// Marshal the request body
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request body"))
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	span.LogFields(log.String("request.body", string(requestBodyBytes)))

	// Create POST request
	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create POST request"))
		return fmt.Errorf("failed to create POST request: %v", err)
	}

	botApiKey := ""
	// prepare bot key
	slackSettings, err := s.postgres.SlackSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get slack settings"))
	}
	if slackSettings == nil {
		span.LogFields(log.String("skip", "slack settings not found"))
		s.log.Warnf("slack settings not found for tenant %s", tenant)
		return nil
	} else {
		botApiKey = slackSettings.AccessToken
	}

	// display last first 8 and last 3 chars
	maskedBotApiKey := ""
	if len(botApiKey) > 11 {
		maskedBotApiKey = botApiKey[:8] + "..." + botApiKey[len(botApiKey)-3:]
	}
	span.LogFields(log.String("bot.api.key", maskedBotApiKey))

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+botApiKey)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform POST request"))
		return fmt.Errorf("failed to perform POST request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		return fmt.Errorf("failed to read response body: %v", err)
	}

	span.LogFields(log.String("response.body", string(responseBody)))

	return nil
}
