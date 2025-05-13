package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	common_enum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/ai/interfaces"
	"github.com/customeros/customeros/packages/server/ai/internal/clients"
	"github.com/customeros/customeros/packages/server/ai/internal/config"
	"github.com/customeros/customeros/packages/server/ai/internal/enum"
	"github.com/customeros/customeros/packages/server/ai/services/common"
)

type GroqService struct {
	config *config.GroqConfig
	client *http.Client
}

func NewGroqService(cfg *config.GroqConfig, warehouseRepos *postgres_repository.WarehouseRepositories) *GroqService {
	DEFAULT_TIMEOUT := 45 * time.Second
	return &GroqService{
		config: cfg,
		client: clients.NewLoggingClient(warehouseRepos.APICallLogRepository, common_enum.VendorGroq, &DEFAULT_TIMEOUT),
	}
}

func (c *GroqService) Ask(ctx context.Context, message *interfaces.AskAIRequest) (*string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "GroqService.Ask")
	defer spans.Finish()
	spans.LogObjectAsJson("request", message)

	if message.Prompt == nil {
		err := errors.New("content (user prompt) cannot be nil")
		spans.TraceError(err)
		return nil, err
	}

	reqBody := buildRequest(message)
	return c.executeWithRetry(ctx, reqBody)
}

func (c *GroqService) createRequest(ctx context.Context, reqBody GroqRequest) (*http.Request, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "GroqService.createRequest")
	defer spans.Finish()

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		err := fmt.Errorf("error marshaling request: %w", err)
		spans.TraceError(err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.Url, bytes.NewReader(jsonBody))
	if err != nil {
		err := fmt.Errorf("error creating request: %w", err)
		spans.TraceError(err)
		return nil, err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.ApiKey)

	return req, nil
}

func (c *GroqService) processResponse(ctx context.Context, resp *http.Response) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "GroqService.processResponse")
	defer spans.Finish()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("error reading response body: %w", err)
		spans.TraceError(err)
		return "", err
	}

	if resp.StatusCode == http.StatusOK {
		return c.handleSuccessResponse(ctx, body)
	}

	return "", c.handleErrorResponse(ctx, resp.StatusCode, body)
}

func (c *GroqService) handleSuccessResponse(ctx context.Context, body []byte) (string, error) {
	spans, _ := telemetry.StartServiceSpan(ctx, "GroqService.handleSuccessResponse")
	defer spans.Finish()

	var data GroqResponse
	if err := json.Unmarshal(body, &data); err != nil {
		err := fmt.Errorf("error decoding response: %w", err)
		spans.TraceError(err)
		return "", err
	}

	if len(data.Choices) < 1 {
		err := fmt.Errorf("empty or invalid response from API")
		spans.TraceError(err)
		return "", err
	}

	resp := data.Choices[0].Message.Content

	if resp == nil || resp == "" {
		err := fmt.Errorf("empty or invalid response from API")
		spans.TraceError(err)
		return "", err
	}

	response, ok := resp.(string)
	if !ok {
		err := fmt.Errorf("empty or invalid response from API")
		spans.TraceError(err)
		return "", err
	}

	return response, nil
}

func (c *GroqService) handleErrorResponse(ctx context.Context, statusCode int, body []byte) error {
	spans, _ := telemetry.StartServiceSpan(ctx, "GroqService.handleErrorResponse")
	defer spans.Finish()
	spans.LogKV("statusCode", statusCode)
	spans.LogKV("body", string(body))

	var errorResponse ErrorResponse
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		err := fmt.Errorf("API request failed with status %d: %s", statusCode, string(body))
		spans.TraceError(err)
		return err
	}
	spans.LogKV("result.errorType", errorResponse.Error.Type)
	spans.LogKV("result.errorMessage", errorResponse.Error.Message)
	return fmt.Errorf("%s: %s", errorResponse.Error.Type, errorResponse.Error.Message)
}

func (c *GroqService) executeWithRetry(ctx context.Context, reqBody GroqRequest) (*string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "GroqService.executeWithRetry")
	defer spans.Finish()

	var lastErr error

	for attempt := 1; attempt <= common.MAX_RETRIES; attempt++ {
		_, jsonErr := json.MarshalIndent(reqBody, "", "  ")
		if jsonErr != nil {
			err := fmt.Errorf("failed to marshal request for logging: %v", jsonErr)
			spans.TraceError(err)
		}
		req, err := c.createRequest(ctx, reqBody)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("error executing request: %w", err)
			if attempt == common.MAX_RETRIES {
				spans.TraceError(err)
				return nil, lastErr
			}
			common.HandleRetry(attempt, err)
			continue
		}
		defer resp.Body.Close()

		response, err := c.processResponse(ctx, resp)
		if err == nil {
			spans.TraceError(err)
			return &response, nil
		}

		lastErr = err
		if !common.ShouldRetry(attempt) {
			break
		}

		common.HandleRetry(attempt, err)
	}

	return nil, lastErr
}

func buildRequest(message *interfaces.AskAIRequest) GroqRequest {
	var outputFormat string
	switch message.OutputFormat {
	case enum.AIOutputJson:
		outputFormat = "json_object"
	default:
		outputFormat = "text"
	}

	// Initialize with default values
	var systemPrompt, prompt string
	var maxCompletionTokens int32
	var temperature float32

	// Safely dereference pointers with nil checks
	if message.SystemPrompt != nil {
		systemPrompt = *message.SystemPrompt
	}

	if message.Prompt != nil {
		prompt = *message.Prompt
	}

	if message.MaxOutputTokens != nil {
		maxCompletionTokens = *message.MaxOutputTokens
	}

	if message.ModelTemperature != nil {
		temperature = *message.ModelTemperature
	}

	req := GroqRequest{
		Model: message.Model.String(),
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
		MaxCompletionTokens: maxCompletionTokens,
		Temperature:         temperature,
		ResponseFormat: ResponseFormat{
			Type: outputFormat,
		},
	}
	return req
}
