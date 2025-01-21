package api_slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
)

type slackService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	slack        interfaces.SlackService
}

func NewSlackService(
	log logger.Logger,
	repositories *repository.Repositories,
	grpcClients *grpc_client.Clients,
	slack interfaces.SlackService,
) cosapi_interfaces.SlackService {
	return &slackService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		slack:        slack,
	}
}

func (s *slackService) GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.GetSlackChannels")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("page", page))
	span.LogFields(log.Object("limit", limit))

	channels, totalCount, err := s.slack.GetPaginatedSlackChannels(ctx, tenant, page, limit)
	if err != nil {
		tracing.TraceErr(span, err)
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

func (s *slackService) GetSlackSettings(ctx context.Context, tenant string) (*cosapi_interfaces.SlackSettingsResponse, error) {
	slackSettings, err := s.repositories.PostgresRepositories.SlackSettingsRepository.Get(ctx, tenant)
	if err != nil {
		return nil, err
	}

	if slackSettings == nil {
		return &cosapi_interfaces.SlackSettingsResponse{
			SlackEnabled: false,
		}, nil
	}

	slackSettingsResponse := cosapi_interfaces.SlackSettingsResponse{
		SlackEnabled: true,
	}

	return &slackSettingsResponse, nil
}

func (s *slackService) getBotToken(ctx context.Context, tenant string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.getBotToken")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	slackSettingsEntity, err := s.repositories.PostgresRepositories.SlackSettingsRepository.Get(ctx, tenant)
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

func (s *slackService) ListSlackChannelsWithBot(ctx context.Context, tenant string) ([]cosapi_interfaces.SlackChannelResponse, error) {
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
		OK       bool                                     `json:"ok"`
		Channels []cosapi_interfaces.SlackChannelResponse `json:"channels"`
		Error    string                                   `json:"error"`
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
