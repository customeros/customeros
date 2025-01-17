package config

type SmartyConfig struct {
	AuthId    string `env:"SMARTY_AUTH_ID"`
	AuthToken string `env:"SMARTY_AUTH_TOKEN"`
}

type IpDataConfig struct {
	ApiUrl             string `env:"IPDATA_API_URL"`
	ApiKey             string `env:"IPDATA_API_KEY"`
	IpDataCacheTtlDays int    `env:"IPDATA_CACHE_TTL_DAYS" envDefault:"90"`
}

type IpHunterConfig struct{}

type EmailConfig struct {
	EmailValidationFromDomain                  string `env:"EMAIL_VALIDATION_FROM_DOMAIN"`
	EmailDomainValidationCacheTtlDays          int    `env:"EMAIL_VALIDATION_DOMAIN_CACHE_TTL_DAYS" envDefault:"90"`
	EmailValidationCacheTtlDays                int    `env:"EMAIL_VALIDATION_CACHE_TTL_DAYS" envDefault:"14"`
	EmailValidationSkipProvidersCommaSeparated string `env:"EMAIL_VALIDATION_SKIP_PROVIDERS" envDefault:""`
}

type MailstackConfig struct {
	SupportedTlds []string `env:"MAILSTACK_SUPPORTED_TLDS" envDefault:"com"`
}

type ScrubbyIoConfig struct {
	ApiUrl       string `env:"SCRUBBY_IO_API_URL" envDefault:"https://api.scrubby.io"`
	ApiKey       string `env:"SCRUBBY_IO_API_KEY"`
	CacheTtlDays int    `env:"SCRUBBY_IO_CACHE_TTL_DAYS" envDefault:"90"`
	CallbackUrl  string `env:"SCRUBBY_IO_CALLBACK_URL"`
}

type TrueInboxConfig struct {
	Enabled      bool   `env:"TRUEINBOX_ENABLED" envDefault:"true"`
	ApiUrl       string `env:"TRUEINBOX_API_URL" envDefault:"https://api.trueinbox.io"`
	ApiKey       string `env:"TRUEINBOX_API_KEY"`
	CacheTtlDays int    `env:"TRUEINBOX_CACHE_TTL_DAYS" envDefault:"30"`
}

type EnrowConfig struct {
	Enabled               bool   `env:"ENROW_ENABLED" envDefault:"true"`
	ApiUrl                string `env:"ENROW_API_URL" envDefault:"https://api.enrow.io"`
	ApiKey                string `env:"ENROW_API_KEY"`
	CacheTtlDays          int    `env:"ENROW_CACHE_TTL_DAYS" envDefault:"14"`
	MaxWaitResultsSeconds int    `env:"ENROW_MAX_WAIT_RESULTS_SECONDS" envDefault:"5"`
	CallbackUrl           string `env:"ENROW_CALLBACK_URL"`
	EnrowCallbackApiKey   string `env:"ENROW_CALLBACK_API_KEY"`
}
