package config

type NylasConfig struct {
	APIKey string `env:"NYLAS_API_KEY"`
	APIUrl string `env:"NYLAS_API_URL"`
}
