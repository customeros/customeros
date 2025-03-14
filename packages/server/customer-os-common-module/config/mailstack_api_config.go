package config

type MailstackApiConfig struct {
	ApiUrl string `env:"MAILSTACK_API_URL"`
	ApiKey string `env:"MAILSTACK_API_KEY"`
	// Deprecated, to be removed
	SupportedTlds []string `env:"MAILSTACK_SUPPORTED_TLD" envDefault:"com"`
}
