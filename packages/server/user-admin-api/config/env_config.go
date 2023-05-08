package config

type Config struct {
	CustomerOS struct {
		CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
	}
	Service struct {
		ServerAddress string `env:"USER_ADMIN_API_SERVER_ADDRESS,required"`
		CorsUrl       string `env:"USER_ADMIN_API_CORS_URL,required"`
		ApiKey        string `env:"USER_ADMIN_API_KEY,required"`
	}
}
