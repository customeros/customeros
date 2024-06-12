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
	Comms struct {
		CommsAPI    string `env:"COMMS_API_PATH,required"`
		CommsAPIKey string `env:"COMMS_MAIL_API_KEY,required"`
	}
	Service struct {
		ServerAddress      string `env:"USER_ADMIN_API_SERVER_ADDRESS,required"`
		CorsUrl            string `env:"USER_ADMIN_API_CORS_URL,required"`
		ApiKey             string `env:"USER_ADMIN_API_KEY,required"`
		ProviderTenantName string `env:"PROVIDER_TENANT_NAME,required"`
		ProviderUsername   string `env:"PROVIDER_USERNAME,required"`
	}
	GoogleOAuth struct {
		ClientId     string `env:"GOOGLE_OAUTH_CLIENT_ID,required"`
		ClientSecret string `env:"GOOGLE_OAUTH_CLIENT_SECRET,required"`
	}
	Slack struct {
		ClientId                         string `env:"SLACK_CLIENT_ID,required"`
		ClientSecret                     string `env:"SLACK_CLIENT_SECRET,required"`
		NotifyNewTenantRegisteredWebhook string `env:"SLACK_NOTIFY_NEW_TENANT_REGISTERED_WEBHOOK,required"`
	}
	GrpcClientConfig config.GrpcClientConfig
	Postgres         config.PostgresConfig
	Neo4j            config.Neo4jConfig
	Jaeger           tracing.JaegerConfig
}
