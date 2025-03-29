package config

type MailSherpaApiConfig struct {
	MailsherpaApiUrl string `env:"MAILSHERPA_API_URL" envDefault:"https://mailsherpa.customeros.cloud"`
	MailsherpaApiKey string `env:"MAILSHERPA_API_KEY" envDefault:""`
}
