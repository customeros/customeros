package ai

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

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

const (
	ApiKeyHeader          = "x-api-key"
	AnthropicApiHeader    = "anthropic-version"
	ContentTypeHeader     = "content-type"
	DefaultApiVersion     = "2023-06-01"
	MaxTokens             = 1024
	MaxRetries            = 4
	DefaultTemperature    = 0.1
	DefaultTimeoutSeconds = 45
)

type ErrorResponse struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

type AnthropicClient struct {
	apiKey string
	apiUrl string
	model  string
	client *http.Client
}

func NewAnthropicClient(cfg *config.AnthropicConfig, model enum.AIModel) *AnthropicClient {
	return &AnthropicClient{
		apiKey: cfg.ApiKey,
		apiUrl: cfg.ApiPath,
		model:  model.String(),
		client: &http.Client{
			Timeout: DefaultTimeoutSeconds * time.Second,
		},
	}
}

func (c *AnthropicClient) Invoke(ctx context.Context, systemPrompt *string, content any) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnthropicClient.Invoke")
	defer span.Finish()

	span.LogKV(
		"systemPrompt", utils.IfNotNilString(systemPrompt),
		"content", content,
	)

	if content == nil {
		err := errors.New("content (user prompt) cannot be nil")
		tracing.TraceErr(span, err)
		return "", err
	}

	reqBody := c.buildRequest(systemPrompt, content)
	return c.executeWithRetry(ctx, reqBody)
}

func (c *AnthropicClient) buildRequest(systemPrompt *string, content any) AnthropicApiRequest {
	// If content is a map, convert it to JSON string
	var processedContent any
	if contentMap, ok := content.(map[string]any); ok {
		jsonBytes, err := json.Marshal(contentMap)
		if err == nil {
			processedContent = string(jsonBytes)
		} else {
			processedContent = content
		}
	} else {
		processedContent = content
	}

	req := AnthropicApiRequest{
		Model:       c.model,
		Messages:    []Message{{Role: "user", Content: []any{processedContent}}},
		MaxTokens:   MaxTokens,
		Temperature: DefaultTemperature,
	}

	if systemPrompt != nil {
		req.System = *systemPrompt
	}

	return req
}

func (c *AnthropicClient) createHttpRequest(ctx context.Context, reqBody AnthropicApiRequest) (*http.Request, error) {
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.apiUrl, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set(ContentTypeHeader, "application/json")
	req.Header.Set(ApiKeyHeader, c.apiKey)
	req.Header.Set(AnthropicApiHeader, DefaultApiVersion)

	return req, nil
}

func (c *AnthropicClient) processResponse(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		return c.handleSuccessResponse(body)
	}

	return "", c.handleErrorResponse(resp.StatusCode, body)
}

func (c *AnthropicClient) handleSuccessResponse(body []byte) (string, error) {
	var data AnthropicApiResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(data.Content) > 0 && data.Content[0].Type == "text" {
		return strings.TrimSpace(data.Content[0].Text), nil
	}
	return "", fmt.Errorf("empty or invalid response from API")
}

func (c *AnthropicClient) handleErrorResponse(statusCode int, body []byte) error {
	var errorResponse ErrorResponse
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		return fmt.Errorf("API request failed with status %d: %s", statusCode, string(body))
	}
	return fmt.Errorf("%s: %s", errorResponse.Error.Type, errorResponse.Error.Message)
}

func (c *AnthropicClient) executeWithRetry(ctx context.Context, reqBody AnthropicApiRequest) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		req, err := c.createHttpRequest(ctx, reqBody)
		if err != nil {
			return "", err
		}

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("error executing request: %w", err)
			if attempt == MaxRetries {
				return "", lastErr
			}
			c.handleRetry(attempt, err)
			continue
		}
		defer resp.Body.Close()

		response, err := c.processResponse(resp)
		if err == nil {
			return response, nil
		}

		lastErr = err
		if !c.shouldRetry(attempt) {
			break
		}

		c.handleRetry(attempt, err)
	}

	logrus.Errorf("Failed after %d attempts: %v", MaxRetries, lastErr)
	return "", lastErr
}

func (c *AnthropicClient) shouldRetry(attempt int) bool {
	if attempt >= MaxRetries {
		return false
	}

	// Add custom retry logic here (e.g., checking for rate limit errors)
	return true
}

func (c *AnthropicClient) handleRetry(attempt int, err error) {
	logrus.Warnf("Request failed (attempt %d/%d): %v", attempt, MaxRetries, err)
	backoff := time.Duration(attempt) * time.Second
	time.Sleep(backoff)
}
