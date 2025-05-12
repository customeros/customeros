package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	common_enum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/sirupsen/logrus"

	"github.com/customeros/customeros/packages/server/ai/interfaces"
	"github.com/customeros/customeros/packages/server/ai/internal/clients"
	"github.com/customeros/customeros/packages/server/ai/internal/config"
	"github.com/customeros/customeros/packages/server/ai/proto/pb"
	"github.com/customeros/customeros/packages/server/ai/services/common"
)

const (
	ApiKeyHeader       = "x-api-key"
	AnthropicApiHeader = "anthropic-version"
	ContentTypeHeader  = "content-type"
	DefaultApiVersion  = "2023-06-01"
)

type AnthropicService struct {
	interfaces.LLMClient
	config *config.AnthropicConfig
	client *http.Client
}

func NewAnthropicService(cfg *config.AnthropicConfig, warehouseRepos *postgres_repository.WarehouseRepositories) *AnthropicService {
	DEFAULT_TIMEOUT := 45 * time.Second
	return &AnthropicService{
		config: cfg,
		client: clients.NewLoggingClient(warehouseRepos.APICallLogRepository, common_enum.VendorAnthropic, &DEFAULT_TIMEOUT),
	}
}

func (c *AnthropicService) Ask(ctx context.Context, message *pb.AskAI) (*string, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "AnthropicService.Ask")
	defer span.Finish()
	span.LogObjectAsJson("request", message)

	if message.Model != pb.AIModel_AI_MODEL_ANTHROPIC_HAIKU && message.Model != pb.AIModel_AI_MODEL_ANTHROPIC_SONNET {
		err := errors.New("model not an anthropic model")
		span.TraceError(err)
		return nil, err
	}
	if *message.Prompt == "" || message.Prompt == nil {
		err := errors.New("content (user prompt) cannot be nil")
		span.TraceError(err)
		return nil, err
	}

	reqBody := buildRequest(message)
	resp, err := c.executeWithRetry(ctx, reqBody)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}
	return &resp, nil
}

func (c *AnthropicService) createRequest(ctx context.Context, reqBody AnthropicApiRequest) (*http.Request, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "AnthropicService.createRequest")
	defer span.Finish()

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		span.TraceError(err)
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.ApiPath, bytes.NewReader(jsonBody))
	if err != nil {
		span.TraceError(err)
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set(ContentTypeHeader, "application/json")
	req.Header.Set(ApiKeyHeader, c.config.ApiKey)
	req.Header.Set(AnthropicApiHeader, DefaultApiVersion)

	return req, nil
}

func (c *AnthropicService) processResponse(ctx context.Context, resp *http.Response) (string, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "AnthropicService.processResponse")
	defer span.Finish()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.TraceError(err)
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		return c.handleSuccessResponse(ctx, body)
	}

	return "", c.handleErrorResponse(ctx, resp.StatusCode, body)
}

func (c *AnthropicService) handleSuccessResponse(ctx context.Context, body []byte) (string, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "AnthropicService.handleSuccessResponse")
	defer span.Finish()

	var data AnthropicApiResponse
	if err := json.Unmarshal(body, &data); err != nil {
		span.TraceError(err)
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(data.Content) > 0 && data.Content[0].Type == "text" {
		return strings.TrimSpace(data.Content[0].Text), nil
	}
	return "", fmt.Errorf("empty or invalid response from API")
}

func (c *AnthropicService) handleErrorResponse(ctx context.Context, statusCode int, body []byte) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "AnthropicService.handleErrorResponse")
	defer span.Finish()

	var errorResponse ErrorResponse
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		span.TraceError(err)
		return fmt.Errorf("API request failed with status %d: %s", statusCode, string(body))
	}
	return fmt.Errorf("%s: %s", errorResponse.Error.Type, errorResponse.Error.Message)
}

func (c *AnthropicService) executeWithRetry(ctx context.Context, reqBody AnthropicApiRequest) (string, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "AnthropicService.executeWithRetry")
	defer span.Finish()

	var lastErr error

	for attempt := 1; attempt <= common.MAX_RETRIES; attempt++ {
		req, err := c.createRequest(ctx, reqBody)
		if err != nil {
			span.TraceError(err)
			return "", err
		}

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("error executing request: %w", err)
			if attempt == common.MAX_RETRIES {
				span.TraceError(err)
				return "", lastErr
			}
			common.HandleRetry(attempt, err)
			continue
		}
		defer resp.Body.Close()

		response, err := c.processResponse(ctx, resp)
		if err == nil {
			return response, nil
		}

		lastErr = err
		if !common.ShouldRetry(attempt) {
			break
		}

		common.HandleRetry(attempt, err)
	}

	logrus.Errorf("Failed after %d attempts: %v", common.MAX_RETRIES, lastErr)
	return "", lastErr
}

func buildRequest(message *pb.AskAI) AnthropicApiRequest {
	req := AnthropicApiRequest{
		Model:       message.Model.String(),
		Messages:    []Message{{Role: "user", Content: message.Prompt}},
		MaxTokens:   *message.MaxOutputTokens,
		Temperature: *message.ModelTemperature,
	}

	if message.SystemPrompt != nil {
		req.System = *message.SystemPrompt
	}

	return req
}
