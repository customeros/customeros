package config

type CloudflareConfig struct {
	Url    string `env:"CLOUDFLARE_URL" envDefault:"https://api.cloudflare.com/client/v4" validate:"required"`
	ApiKey string `env:"CLOUDFLARE_API_KEY" `
	Email  string `env:"CLOUDFLARE_API_EMAIL"`
}

type OpenSRSConfig struct {
	Url      string `env:"OPENSRS_URL" envDefault:"https://admin.a.hostedemail.com"`
	ApiKey   string `env:"OPENSRS_API_KEY"`
	Username string `env:"OPENSRS_API_USERNAME"`
}

type PostmarkConfig struct {
	Url                         string `env:"POSTMARK_URL" envDefault:"https://api.postmarkapp.com"`
	AccountApiKey               string `env:"POSTMARK_ACCOUNT_API_KEY"`
	DefaultInboundStreamWebhook string `env:"POSTMARK_DEFAULT_INBOUND_STREAM_WEBHOOK"`
}

type StripeConfig struct {
	ApiKey string `env:"STRIPE_API_KEY" envDefault:"N/A"`
}
