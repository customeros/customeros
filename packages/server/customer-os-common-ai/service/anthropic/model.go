package anthropic

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/config"
)

type AnthropicModel struct {
	Client *AnthropicClient
}

func NewModel(apiKey, apiPath, model string) *AnthropicModel {
	cfg := &config.AiModelConfigAnthropic{
		ApiKey:  apiKey,
		ApiPath: apiPath,
		Model:   model,
	}
	return &AnthropicModel{
		Client: NewAnthropicClient(cfg),
	}
}

func (m *AnthropicModel) Inference(ctx context.Context, input string) (string, error) {
	return m.Client.Invoke(ctx, input)
}
