package config

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"

type Config struct {
	Logger     logger.Config
	CustomerOS struct {
		CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
	}
	Service struct {
		ServerAddress string `env:"USER_ADMIN_API_SERVER_ADDRESS,required"`
		CorsUrl       string `env:"USER_ADMIN_API_CORS_URL,required"`
		ApiKey        string `env:"USER_ADMIN_API_KEY,required"`
	}
	GoogleOAuth struct {
		ClientId     string `env:"GOOGLE_OAUTH_CLIENT_ID,required"`
		ClientSecret string `env:"GOOGLE_OAUTH_CLIENT_SECRET,required"`
	}
	Postgres struct {
		Host            string `env:"POSTGRES_HOST,required"`
		Port            string `env:"POSTGRES_PORT,required"`
		User            string `env:"POSTGRES_USER,required,unset"`
		Db              string `env:"POSTGRES_DB,required"`
		Password        string `env:"POSTGRES_PASSWORD,required,unset"`
		MaxConn         int    `env:"POSTGRES_DB_MAX_CONN"`
		MaxIdleConn     int    `env:"POSTGRES_DB_MAX_IDLE_CONN"`
		ConnMaxLifetime int    `env:"POSTGRES_DB_CONN_MAX_LIFETIME"`
		LogLevel        string `env:"POSTGRES_LOG_LEVEL" envDefault:"WARN"`
	}
}
