package config

import (
	authConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/config"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	Service struct {
		CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
		ServerAddress    string `env:"COMMS_API_SERVER_ADDRESS,required"`
		PublicPath       string `env:"COMMS_API_PUBLIC_PATH,required"`
	}
	Mail struct {
		ApiKey string `env:"COMMS_API_MAIL_API_KEY,required"`
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
	CalCom struct {
		CalComWebhookSecret string `env:"CALCOM_SECRET,required"`
	}
	Redis struct {
		Host   string `env:"REDIS_HOST,required"`
		Scheme string `env:"REDIS_SCHEME,required"envDefault:"rediss"`
	}
	FileStoreApiConfig fsc.FileStoreApiConfig
	AuthConfig         authConfig.Config
	Jaeger             tracing.JaegerConfig
	Logger             logger.Config
}
