package config

type CustomerOsApiConfig struct {
	ApiUrl string `env:"CUSTOMER_OS_API_URL" envDefault:"https://api.customeros.ai"`
	ApiKey string `env:"CUSTOMER_OS_API_KEY"`
}
