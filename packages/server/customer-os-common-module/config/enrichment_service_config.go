package config

type ScrapinConfig struct {
	Url     string `env:"SCRAPIN_API_URL" envDefault:"https://api.scrapin.io" required:"true"`
	ApiKey  string `env:"SCRAPIN_API_KEY" required:"true"`
	TtlDays int    `env:"SCRAPIN_TTL_DAYS" envDefault:"90" required:"true"`
}

type BrandfetchConfig struct {
	Url     string `env:"BRANDFETCH_API_URL"`
	Limit   int    `env:"BRANDFETCH_LIMIT" envDefault:"250"`
	TtlDays int    `env:"BRANDFETCH_TTL_DAYS" envDefault:"180" required:"true"`
}

type BetterContactConfig struct {
	Url         string `env:"BETTER_CONTACT_API_URL" required:"true"`
	ApiKey      string `env:"BETTER_CONTACT_API_KEY" required:"true"`
	CallbackUrl string `env:"BETTER_CONTACT_CALLBACK_URL" required:"true"`
}

type SnitcherConfig struct {
	Url    string `env:"SNITCHER_API_URL" required:"true" envDefault:"https://app.snitcher.com/api"`
	ApiKey string `env:"SNITCHER_API_KEY" required:"true"`
}
