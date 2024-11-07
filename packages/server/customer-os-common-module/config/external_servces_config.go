package config

type OpenSRSConfig struct {
	Url      string `env:"OPENSRS_URL" envDefault:"https://admin.a.hostedemail.com"`
	ApiKey   string `env:"OPENSRS_API_KEY"`
	Username string `env:"OPENSRS_API_USERNAME"`
}
