package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type DeepseekClient struct {
	apiUrl string
	apiKey string
	client *http.Client
}

func NewDeepseekClient(cfg *config.DeepseekConfig) *DeepseekClient {
	return &DeepseekClient{
		apiKey: cfg.ApiKey,
		apiUrl: cfg.Url,
		client: &http.Client{
			Timeout: DefaultTimeoutSeconds * time.Second,
		},
	}
}

func (c *DeepseekClient) AskDeepseek(ctx context.Context, request interfaces.AskAIRequest) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DeepseekClient.AskDeepseek")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "request", request)

	if request.Model != enum.AIModelDeepseekChat {
		err := errors.New("model not a deepseek model")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if request.Prompt == nil {
		err := errors.New("content cannot be nil")
		tracing.TraceErr(span, err)
		return nil, err
	}

	req := c.buildRequest(request)
	return c.executeWithRetry(ctx, req)
}

func (c *DeepseekClient) buildRequest(request interfaces.AskAIRequest) *DeepseekRequest {
	sysPromptStr := "You are a helpful assistant."
	if request.SystemPrompt != nil {
		sysPromptStr = *request.SystemPrompt
	}

	return &DeepseekRequest{
		Model: request.Model.String(),
		Messages: []Message{
			{Role: "system", Content: sysPromptStr},
			{Role: "user", Content: *request.Prompt},
		},
		Stream: false,
	}
}

func (c *DeepseekClient) executeWithRetry(ctx context.Context, req *DeepseekRequest) (*string, error) {
	var lastErr error

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		resp, err := c.ChatCompletions(ctx, req)
		if err == nil && len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
			answer := resp.Choices[0].Message.Content
			return &answer, nil
		}

		lastErr = err
		if !c.shouldRetry(attempt) {
			break
		}

		c.handleRetry(attempt, err)
	}

	logrus.Errorf("Failed after %d attempts: %v", MaxRetries, lastErr)
	return nil, lastErr
}

func (c *DeepseekClient) ChatCompletions(ctx context.Context, req *DeepseekRequest) (*DeepseekResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", c.apiUrl+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	return c.processResponse(resp)
}

func (c *DeepseekClient) processResponse(resp *http.Response) (*DeepseekResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response DeepseekResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &response, nil
}

func (c *DeepseekClient) shouldRetry(attempt int) bool {
	return attempt < MaxRetries
}

func (c *DeepseekClient) handleRetry(attempt int, err error) {
	logrus.Warnf("Request failed (attempt %d/%d): %v", attempt, MaxRetries, err)
	backoff := time.Duration(attempt) * time.Second
	time.Sleep(backoff)
}
