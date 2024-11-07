package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	Logger     logger.Config
	CustomerOS struct {
		CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
	}
	Service struct {
		Port               string `env:"PORT,required"`
		PublicPath         string `env:"USER_ADMIN_API_PUBLIC_PATH,required"`
		CorsUrl            string `env:"CORS_URL,required"`
		ProviderTenantName string `env:"PROVIDER_TENANT_NAME,required"`
		ProviderUsername   string `env:"PROVIDER_USERNAME,required"`
	}
	GoogleOAuth config.GoogleOAuthConfig
	Slack       struct {
		ClientId                         string `env:"SLACK_CLIENT_ID,required"`
		ClientSecret                     string `env:"SLACK_CLIENT_SECRET,required"`
		NotifyNewTenantRegisteredWebhook string `env:"SLACK_NOTIFY_NEW_TENANT_REGISTERED_WEBHOOK,required"`
	}
	GrpcClientConfig config.GrpcClientConfig
	Postgres         config.PostgresConfig
	Neo4j            config.Neo4jConfig
	Jaeger           tracing.JaegerConfig
	RabbitMQConfig   config.RabbitMQConfig
	OpenSRSConfig    config.OpenSRSConfig
}
