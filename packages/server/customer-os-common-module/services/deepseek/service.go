package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type DeepseekService struct {
	config *config.DeepseekConfig
}

func NewDeepseekService(config *config.DeepseekConfig) *DeepseekService {
	return &DeepseekService{
		config: config,
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepseekRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type DeepseekResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	LogProbs     any     `json:"logprobs"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens          int                 `json:"prompt_tokens"`
	CompletionTokens      int                 `json:"completion_tokens"`
	TotalTokens           int                 `json:"total_tokens"`
	PromptTokensDetails   PromptTokensDetails `json:"prompt_tokens_details"`
	PromptCacheHitTokens  int                 `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens int                 `json:"prompt_cache_miss_tokens"`
}

type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

func (s *DeepseekService) AskDeepseek(ctx context.Context, model enum.AIModel, systemPrompt *string, prompt string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DeepseekService.AskDeepseek")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("model", model, "systemPrompt", utils.IfNotNilString(systemPrompt))
	tracing.LogObjectAsJson(span, "prompt", prompt)

	if model != enum.AIModelDeepseekChat {
		err := errors.New("model not a deepseek model")
		tracing.TraceErr(span, err)
		return nil, err
	}

	sysPromptStr := "You are a helpful assistant."
	if systemPrompt != nil {
		sysPromptStr = *systemPrompt
	}

	req := DeepseekRequest{
		Model: model.String(),
		Messages: []Message{
			{Role: "system", Content: sysPromptStr},
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", s.config.Url+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+s.config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response DeepseekResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("no response content from AI")
	}

	answer := response.Choices[0].Message.Content
	return &answer, nil
}
