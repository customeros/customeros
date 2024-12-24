package service

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/service/anthropic"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/service/openai"
)

// AiModel interface defines the common behavior for all AI models
type AiModel interface {
	Inference(ctx context.Context, input string) (string, error)
}

type AiModelType string

const (
	AnthropicModelType AiModelType = "anthropic"
	OpenAiModelType    AiModelType = "openai"
)

// NewAiModel is a factory function that returns the appropriate AI model based on type
func NewAiModel(modelType AiModelType, cfg config.Config) AiModel {
	switch modelType {
	case AnthropicModelType:
		return NewAnthropicModel(cfg.Anthropic.ApiKey, cfg.Anthropic.ApiPath, cfg.Anthropic.Model)
	case OpenAiModelType:
		return NewOpenAiModel(cfg.OpenAi.ApiKey, cfg.OpenAi.Organization, cfg.OpenAi.Model)
	default:
		return nil
	}
}

/////////////////////// OpenAI ///////////////////////

// NewOpenAiModel creates a new OpenAI model instance
func NewOpenAiModel(apiKey, organization, model string) AiModel {
	return openai.NewModel(apiKey, organization, model)
}

//////////////// Anthropic ///////////////////////

// NewAnthropicModel creates a new Anthropic model instance
// This now uses your new Anthropic client implementation
func NewAnthropicModel(apiKey, apiPath, model string) AiModel {
	return anthropic.NewModel(apiKey, apiPath, model)
}
