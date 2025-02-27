package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type GroqClient struct {
	apiKey string
	apiUrl string
	client *http.Client
}

func NewGroqClient(cfg *config.GroqConfig) *GroqClient {
	return &GroqClient{
		apiKey: cfg.ApiKey,
		apiUrl: cfg.Url,
		client: &http.Client{
			Timeout: DefaultTimeoutSeconds * time.Second,
		},
	}
}

func (c *GroqClient) Invoke(ctx context.Context, request interfaces.AskAIRequest) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GroqClient.Invoke")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "request", request)

	if request.Prompt == nil {
		err := errors.New("content (user prompt) cannot be nil")
		tracing.TraceErr(span, err)
		return "", err
	}

	reqBody := c.buildRequest(request)
	return c.executeWithRetry(ctx, reqBody)
}

func (c *GroqClient) buildRequest(request interfaces.AskAIRequest) GroqRequest {
	var outputFormat string
	switch request.OutputFormat {
	case enum.AIOutputText:
		outputFormat = "text"
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
	if request.SystemPrompt != nil {
		systemPrompt = *request.SystemPrompt
	}

	if request.Prompt != nil {
		prompt = *request.Prompt
	}

	if request.MaxOutputTokens != nil {
		maxCompletionTokens = *request.MaxOutputTokens
	}

	if request.ModelTemperature != nil {
		temperature = *request.ModelTemperature
	}

	req := GroqRequest{
		Model: request.Model.String(),
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

func (c *GroqClient) createHttpRequest(ctx context.Context, reqBody GroqRequest) (*http.Request, error) {
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.apiUrl, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set(ContentTypeHeader, "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	return req, nil
}

func (c *GroqClient) processResponse(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		return c.handleSuccessResponse(body)
	}

	return "", c.handleErrorResponse(resp.StatusCode, body)
}

func (c *GroqClient) handleSuccessResponse(body []byte) (string, error) {
	var data GroqResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(data.Choices) < 1 {
		return "", fmt.Errorf("empty or invalid response from API")
	}

	resp := data.Choices[0].Message.Content

	if resp == nil || resp == "" {
		return "", fmt.Errorf("empty or invalid response from API")
	}

	response, ok := resp.(string)
	if !ok {
		return "", fmt.Errorf("empty or invalid response from API")
	}

	return response, nil
}

func (c *GroqClient) handleErrorResponse(statusCode int, body []byte) error {
	var errorResponse ErrorResponse
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		return fmt.Errorf("API request failed with status %d: %s", statusCode, string(body))
	}
	return fmt.Errorf("%s: %s", errorResponse.Error.Type, errorResponse.Error.Message)
}

func (c *GroqClient) executeWithRetry(ctx context.Context, reqBody GroqRequest) (string, error) {
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

func (c *GroqClient) shouldRetry(attempt int) bool {
	if attempt >= MaxRetries {
		return false
	}

	// Add custom retry logic here (e.g., checking for rate limit errors)
	return true
}

func (c *GroqClient) handleRetry(attempt int, err error) {
	logrus.Warnf("Request failed (attempt %d/%d): %v", attempt, MaxRetries, err)
	backoff := time.Duration(attempt) * time.Second
	time.Sleep(backoff)
}
