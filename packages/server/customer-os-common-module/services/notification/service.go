package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type notificationService struct {
	log                  logger.Logger
	postgresRepositories *repository.Repositories
}

func NewNotificationService(log logger.Logger, postgresRepo *repository.Repositories) interfaces.NotificationService {
	return &notificationService{
		log:                  log,
		postgresRepositories: postgresRepo,
	}
}

func (s *notificationService) NotifySlackChannel(ctx context.Context, tenant, channelID string, message *string) error {
	span, ctx := tracing.StartTracerSpan(ctx, "NotificationService.NotifySlackChannel")
	defer span.Finish()

	err := s.sendSlackMessage(ctx, tenant, channelID, *message)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *notificationService) sendSlackMessage(ctx context.Context, tenant, channel, blocks string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "NotificationService.sendSlackMessage")
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
	slackSettings, err := s.postgresRepositories.SlackSettingsRepository.Get(ctx, tenant)
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
