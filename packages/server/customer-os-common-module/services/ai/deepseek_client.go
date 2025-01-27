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
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type DeepseekClient struct {
	apiUrl string
	apiKey string
	model  string
	client *http.Client
}

func NewDeepseekClient(cfg *config.DeepseekConfig, model enum.AIModel) *DeepseekClient {
	return &DeepseekClient{
		apiKey: cfg.ApiKey,
		apiUrl: cfg.Url,
		model:  model.String(),
		client: &http.Client{
			Timeout: DefaultTimeoutSeconds * time.Second,
		},
	}
}

func (c *DeepseekClient) AskDeepseek(ctx context.Context, systemPrompt string, content any) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DeepseekClient.AskDeepseek")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("model", c.model, "systemPrompt", utils.IfNotNilString(systemPrompt))
	tracing.LogObjectAsJson(span, "content", content)

	if c.model != enum.AIModelDeepseekChat.String() {
		err := errors.New("model not a deepseek model")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if content == nil {
		err := errors.New("content cannot be nil")
		tracing.TraceErr(span, err)
		return nil, err
	}

	req := c.buildRequest(systemPrompt, content)
	return c.executeWithRetry(ctx, req)
}

func (c *DeepseekClient) buildRequest(systemPrompt string, content any) *DeepseekRequest {
	sysPromptStr := "You are a helpful assistant."
	if systemPrompt != "" {
		sysPromptStr = systemPrompt
	}

	var promptContent string
	switch p := content.(type) {
	case string:
		promptContent = p
	default:
		jsonBytes, err := json.Marshal(p)
		if err != nil {
			return nil
		}
		promptContent = string(jsonBytes)
	}

	return &DeepseekRequest{
		Model: c.model,
		Messages: []Message{
			{Role: "system", Content: sysPromptStr},
			{Role: "user", Content: promptContent},
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
