package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	common_enum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/customeros/customeros/packages/server/ai/interfaces"
	"github.com/customeros/customeros/packages/server/ai/internal/clients"
	"github.com/customeros/customeros/packages/server/ai/internal/config"
	"github.com/customeros/customeros/packages/server/ai/internal/enum"
	"github.com/customeros/customeros/packages/server/ai/services/common"
)

const (
	MAX_TOKENS          = 1024
	DEFAULT_TEMPERATURE = 0.1
)

type DeepseekService struct {
	config *config.DeepseekConfig
	client *http.Client
}

func NewDeepseekService(cfg *config.DeepseekConfig, warehouseRepos *postgres_repository.WarehouseRepositories) *DeepseekService {
	DEFAULT_TIMEOUT := 45 * time.Second

	return &DeepseekService{
		config: cfg,
		client: clients.NewLoggingClient(warehouseRepos.APICallLogRepository, common_enum.VendorDeepseek, &DEFAULT_TIMEOUT),
	}
}

func (c *DeepseekService) Ask(ctx context.Context, message interfaces.AskAIRequest) (*string, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "DeepseekService.Ask")
	defer span.Finish()
	span.LogObjectAsJson("request", message)

	if message.Model != enum.AIModelDeepseekChat {
		err := errors.New("model not a deepseek model")
		span.TraceError(err)
		return nil, err
	}

	if message.Prompt == nil {
		err := errors.New("prompt cannot be nil")
		span.TraceError(err)
		return nil, err
	}

	req := buildRequest(message)
	return c.executeWithRetry(ctx, req)
}

func (c *DeepseekService) executeWithRetry(ctx context.Context, req *DeepseekRequest) (*string, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "DeepseekClient.executeWithRetry")
	defer span.Finish()

	var lastErr error

	for attempt := 1; attempt <= common.MAX_RETRIES; attempt++ {
		resp, err := c.chatCompletions(ctx, req)
		if err == nil && len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
			content := resp.Choices[0].Message.Content
			contentStr := content.(string)
			return &contentStr, nil
		}

		lastErr = err
		if !common.ShouldRetry(attempt) {
			break
		}

		common.HandleRetry(attempt, err)
	}

	logrus.Errorf("Failed after %d attempts: %v", common.MAX_RETRIES, lastErr)
	return nil, lastErr
}

func (c *DeepseekService) chatCompletions(ctx context.Context, req *DeepseekRequest) (*DeepseekResponse, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "DeepseekService.chatCompletions")
	defer span.Finish()

	jsonData, err := json.Marshal(req)
	if err != nil {
		span.TraceError(err)
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", c.config.Url+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		span.TraceError(err)
		return nil, fmt.Errorf("create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.config.ApiKey)

	resp, err := c.client.Do(request)
	if err != nil {
		span.TraceError(err)
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	return c.processResponse(ctx, resp)
}

func (c *DeepseekService) processResponse(ctx context.Context, resp *http.Response) (*DeepseekResponse, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "DeepseekService.processResponse")
	defer span.Finish()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.TraceError(err)
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		span.TraceError(err)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response DeepseekResponse
	if err := json.Unmarshal(body, &response); err != nil {
		span.TraceError(err)
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &response, nil
}

func buildRequest(message interfaces.AskAIRequest) *DeepseekRequest {
	sysPromptStr := "You are a helpful assistant."
	if message.SystemPrompt != nil {
		sysPromptStr = *message.SystemPrompt
	}

	return &DeepseekRequest{
		Model: message.Model.String(),
		Messages: []Message{
			{Role: "system", Content: sysPromptStr},
			{Role: "user", Content: *message.Prompt},
		},
		Stream: false,
	}
}
