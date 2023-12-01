package anthropic

import (
	"context"
)

func NewModel(apiKey, apiPath string) *AnthropicModel {
	cfg := &AiModelConfigAnthropic{
		ApiKey:  apiKey,
		ApiPath: apiPath,
	}
	return &AnthropicModel{
		Client: NewAnthropicClient(cfg),
	}
}

type AnthropicModel struct {
	Client *AnthropicClient
}

func (m *AnthropicModel) Inference(ctx context.Context, input string) (string, error) {
	return InvokeAnthropic(ctx, m.Client.cfg, input)
}
