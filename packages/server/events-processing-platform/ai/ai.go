package ai

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type AiModel interface {
	Inference(ctx context.Context, input string) (string, error)
}

type AiModelType string

const (
	AnthropicModelType AiModelType = "anthropic"
	OpenAiModelType    AiModelType = "openai"
)

func NewAiModel(modelType AiModelType, apiKey, apiPath, organization, model string, logger logger.Logger) AiModel {
	switch modelType {
	case AnthropicModelType:
		return NewAnthropicModel(apiKey, apiPath, organization, logger)
	case OpenAiModelType:
		return NewOpenAiModel(apiKey, organization, model)
	default:
		return nil
	}
}

type any interface{}

type AiModelConfig struct {
	Type AiModelType `json:"type"`
}

type AiModelConfigAnthropic struct {
	AiModelConfig
	ApiKey       string `json:"apiKey"`
	Organization string `json:"organization"`
	ApiPath      string `json:"apiPath"`
}

type AiModelConfigOpenAi struct {
	AiModelConfig
	ApiKey       string `json:"apiKey"`
	Organization string `json:"organization"`
	Model        string `json:"model"`
}

// TODO: make submodules for each ai provider/client
/////////////////////// OpenAI ///////////////////////

func NewOpenAiModel(apiKey, organization, model string) AiModel {
	// TODO: use openai config object plus add logger to openai client
	return &OpenAiModel{
		Client: NewOpenAiClient(apiKey, organization, model),
	}
}

type OpenAiModel struct {
	Client *OpenAiClient
}

func (m *OpenAiModel) Inference(ctx context.Context, input string) (string, error) {
	request := &CreateChatCompletionsRequest{
		Model: m.Client.model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant designed to output JSON.",
			},
			{
				Role:    "user",
				Content: input,
			},
		},
		Temperature:    0.7,
		Seed:           42,                                  // always use same seed to increase likelihood of consistent results
		ResponseFormat: ResponseFormat{Type: "json_object"}, // https://platform.openai.com/docs/guides/text-generation/json-mode
	}

	response, err := m.Client.CreateChatCompletions(ctx, request)
	if err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content, nil
}

//////////////// Anthropic ///////////////////////

func NewAnthropicModel(apiKey, apiPath, organization string, logger logger.Logger) AiModel {
	cfg := &AiModelConfigAnthropic{
		AiModelConfig: AiModelConfig{
			Type: AnthropicModelType,
		},
		ApiKey:       apiKey,
		Organization: organization,
		ApiPath:      apiPath,
	}
	return &AnthropicModel{
		Client: NewAnthropicClient(cfg, logger),
	}
}

type AnthropicModel struct {
	Client *AnthropicClient
}

func (m *AnthropicModel) Inference(ctx context.Context, input string) (string, error) {
	return InvokeAnthropic(ctx, m.Client.cfg, m.Client.logger, input)
}
