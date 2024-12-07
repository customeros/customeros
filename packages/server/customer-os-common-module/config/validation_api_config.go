package config

type ValidationAPIConfig struct {
	Url    string `env:"VALIDATION_API_URL"`
	ApiKey string `env:"VALIDATION_API_KEY"`
}
