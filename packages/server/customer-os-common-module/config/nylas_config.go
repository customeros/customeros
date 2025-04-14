package config

type NylasConfig struct {
	ClientID     string `env:"NYLAS_CLIENT_ID"`
	ClientSecret string `env:"NYLAS_CLIENT_SECRET"`
	APIKey       string `env:"NYLAS_API_KEY"`
	APIURL       string `env:"NYLAS_API_URL" envDefault:"https://api.nylas.com"`
	WebhookURL   string `env:"NYLAS_WEBHOOK_URL"`
}
