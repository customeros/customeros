package service

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/service/anthropic"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/service/openai"
)

// AiModel interface defines the common behavior for all AI models
type AIModel interface {
	Inference(ctx context.Context, input string) (string, error)
}

// NewAiModel is a factory function that returns the appropriate AI model based on type
func NewAIModel(model enum.AIModel, apiKey string) AIModel {
	switch model {
	case enum.AIModelAnthropicHaiku, enum.AIModelAnthropicSonnet:
		cfg := config.NewAnthropicConfig(model, apiKey)
		return NewAnthropicModel(cfg)
	// case OpenAiModelType:
	// 	return NewOpenAiModel(cfg.OpenAi.ApiKey, cfg.OpenAi.Organization, cfg.OpenAi.Model)
	default:
		return nil
	}
}

/////////////////////// OpenAI ///////////////////////

// NewOpenAiModel creates a new OpenAI model instance
func NewOpenAiModel(apiKey, organization, model string) AIModel {
	return openai.NewModel(apiKey, organization, model)
}

//////////////// Anthropic ///////////////////////

// NewAnthropicModel creates a new Anthropic model instance
// This now uses your new Anthropic client implementation
func NewAnthropicModel(cfg *config.AiModelConfigAnthropic) AIModel {
	return anthropic.NewModel(cfg.ApiKey, cfg.ApiPath, cfg.Model)
}
