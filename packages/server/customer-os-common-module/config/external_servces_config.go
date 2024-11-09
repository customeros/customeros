package config

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
