package config

type AiModelConfigAnthropic struct {
	ApiKey  string `json:"apiKey"`
	ApiPath string `json:"apiPath"`
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
