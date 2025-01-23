package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type slackService struct {
	log                  logger.Logger
	postgresRepositories *postgres_repository.Repositories
}

func NewSlackService(log logger.Logger, postgres *postgres_repository.Repositories) interfaces.SlackService {
	return &slackService{
		log:                  log,
		postgresRepositories: postgres,
	}
}

func (s *slackService) GetSlackChannels(ctx context.Context, tenant string) ([]*postgresEntity.SlackChannel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.GetSlackChannels")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, "service")

	nodes, err := s.postgresRepositories.SlackChannelRepository.GetSlackChannels(ctx, tenant)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (s *slackService) GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.GetSlackChannels")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, "service")
	span.LogFields(log.Object("page", page))
	span.LogFields(log.Object("limit", limit))

	channels, totalCount, err := s.postgresRepositories.SlackChannelRepository.GetPaginatedSlackChannels(ctx, tenant, page, limit)
	if err != nil {
		return nil, err
	}

	paginatedResult := utils.Pagination{
		Limit:     page,
		Page:      limit,
		TotalRows: totalCount,
		Rows:      channels,
	}

	return &paginatedResult, nil
}

func (s *slackService) StoreSlackChannel(ctx context.Context, tenant, source, channelId, channelName string, organizationId *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.StoreSlackChannel")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, "service")
	span.LogFields(log.String("channelId", channelId))
	span.LogFields(log.String("channelName", channelName))
	span.LogFields(log.String("organizationId", utils.IfNotNilString(organizationId)))

	existing, err := s.postgresRepositories.SlackChannelRepository.GetSlackChannel(ctx, tenant, channelId)
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
		return s.postgresRepositories.SlackChannelRepository.CreateSlackChannel(ctx, &slackChannel)
	}
	if existing != nil {
		if organizationId != nil {
			return s.postgresRepositories.SlackChannelRepository.UpdateSlackChannelOrganization(ctx, existing.ID, *organizationId)
		} else if channelName != "" {
			return s.postgresRepositories.SlackChannelRepository.UpdateSlackChannelName(ctx, existing.ID, channelName)
		}
	}

	return nil
}

func (s *slackService) SendMessageFromBot(ctx context.Context, channel, blocks string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "NotificationService.sendSlackMessage")
	defer span.Finish()
	span.LogFields(log.String("channel", channel))

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set on context")
		tracing.TraceErr(span, err)
		return err
	}
	span.SetTag(tracing.SpanTagTenant, tenant)

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

func (s *slackService) GetSlackSettings(ctx context.Context, tenant string) (*interfaces.SlackSettingsResponse, error) {
	slackSettings, err := s.postgresRepositories.SlackSettingsRepository.Get(ctx, tenant)
	if err != nil {
		return nil, err
	}

	if slackSettings == nil {
		return &interfaces.SlackSettingsResponse{
			SlackEnabled: false,
		}, nil
	}

	slackSettingsResponse := interfaces.SlackSettingsResponse{
		SlackEnabled: true,
	}

	return &slackSettingsResponse, nil
}

func (s *slackService) getBotToken(ctx context.Context, tenant string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.getBotToken")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	slackSettingsEntity, err := s.postgresRepositories.SlackSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if slackSettingsEntity == nil {
		err := fmt.Errorf("slack settings not found for tenant %s", tenant)
		tracing.TraceErr(span, err)
		return "", err
	}

	return slackSettingsEntity.AccessToken, nil
}

func (s *slackService) ListSlackChannelsWithBot(ctx context.Context, tenant string) ([]interfaces.SlackChannelResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.ListSlackChannelsWithBot")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	token, err := s.getBotToken(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Slack API URL for listing channels
	url := "https://slack.com/api/conversations.list"

	// Build the HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Add required headers
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var result struct {
		OK       bool                              `json:"ok"`
		Channels []interfaces.SlackChannelResponse `json:"channels"`
		Error    string                            `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Handle errors from Slack API
	if !result.OK {
		err := fmt.Errorf("slack API error: %s", result.Error)
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result.Channels, nil
}

func (s *slackService) JoinSlackChannelsWithBot(ctx context.Context, tenant, channelId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.JoinSlackChannelsWithBot")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("channelId", channelId))

	token, err := s.getBotToken(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Slack API URL for joining a channel
	url := "https://slack.com/api/conversations.join"

	// Prepare the request payload
	payload := map[string]string{"channel": channelId}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Build the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Add required headers
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	defer resp.Body.Close()

	// Parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	var result struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Handle errors from Slack API
	if !result.OK {
		err := fmt.Errorf("slack API error: %s", result.Error)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *slackService) LeaveSlackChannelsWithBot(ctx context.Context, tenant, channelId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.LeaveSlackChannelsWithBot")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("channelId", channelId))

	token, err := s.getBotToken(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Slack API URL for leaving a channel
	url := "https://slack.com/api/conversations.leave"

	// Prepare the request payload
	payload := map[string]string{"channel": channelId}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Build the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Add required headers
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	var result struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Handle errors from Slack API
	if !result.OK {
		err := fmt.Errorf("slack API error: %s", result.Error)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
