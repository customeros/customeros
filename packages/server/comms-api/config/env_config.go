package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	Service struct {
		CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
		Port             string `env:"PORT,required"`
		PublicPath       string `env:"COMMS_API_PUBLIC_PATH,required"`
	}
	Mail struct {
		ApiKey string `env:"COMMS_API_MAIL_API_KEY,required"`
	}
	Postgres    config.PostgresConfig
	Neo4jConfig config.Neo4jConfig
	CalCom      struct {
		CalComWebhookSecret string `env:"CALCOM_SECRET,required"`
	}
	Redis struct {
		Host   string `env:"REDIS_HOST,required"`
		Scheme string `env:"REDIS_SCHEME,required"envDefault:"rediss"`
	}
	FileStoreApiConfig fsc.FileStoreApiConfig
	AuthConfig         config.GoogleOAuthConfig
	Jaeger             tracing.JaegerConfig
	Logger             logger.Config
}
