package config

type NovuConfig struct {
	FronteraUrl string `env:"FRONTERA_URL" envDefault:"N/A"`
	ApiKey      string `env:"NOVU_API_KEY" envDefault:"N/A"`
}
