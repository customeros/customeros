package service

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/service/anthropic"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/service/openai"
)

type AiModel interface {
	Inference(ctx context.Context, input string) (string, error)
}

type AiModelType string

const (
	AnthropicModelType AiModelType = "anthropic"
	OpenAiModelType    AiModelType = "openai"
)

func NewAiModel(modelType AiModelType, cfg config.Config) AiModel {
	switch modelType {
	case AnthropicModelType:
		return NewAnthropicModel(cfg.Anthropic.ApiKey, cfg.Anthropic.ApiPath)
	case OpenAiModelType:
		return NewOpenAiModel(cfg.OpenAi.ApiKey, cfg.OpenAi.Organization, cfg.OpenAi.Model)
	default:
		return nil
	}
}

/////////////////////// OpenAI ///////////////////////

func NewOpenAiModel(apiKey, organization, model string) AiModel {
	return openai.NewModel(apiKey, organization, model)
}

//////////////// Anthropic ///////////////////////

func NewAnthropicModel(apiKey, apiPath string) AiModel {
	return anthropic.NewModel(apiKey, apiPath)
}
