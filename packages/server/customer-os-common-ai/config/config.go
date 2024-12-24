package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/enum"
)

type AiModelConfigAnthropic struct {
	ApiKey  string `json:"apiKey" env:"ANTHROPIC_API_KEY"`
	ApiPath string `json:"apiPath" env:"ANTHROPIC_API_PATH"`
	Model   string `json:"model"`
}

type AiModelConfigOpenAi struct {
	ApiKey       string `json:"apiKey"`
	Organization string `json:"organization"`
	Model        string `json:"model"`
}

type Config struct {
	OpenAi    AiModelConfigOpenAi
	Anthropic AiModelConfigAnthropic
}

func NewAnthropicConfig(model enum.AIModel, apiKey string) *AiModelConfigAnthropic {
	return &AiModelConfigAnthropic{
		ApiKey:  apiKey,
		ApiPath: "https://api.anthropic.com/v1/messages",
		Model:   model.String(),
	}
}
