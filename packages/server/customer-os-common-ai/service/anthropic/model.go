package anthropic

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/config"
)

func NewModel(apiKey, apiPath string) *AnthropicModel {
	cfg := &config.AiModelConfigAnthropic{
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
