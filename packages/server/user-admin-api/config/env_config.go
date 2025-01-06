package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

// Config is the main configuration structure
type Config struct {
	Logger           logger.Config
	CustomerOS       CustomerOSConfig
	Service          ServiceConfig
	GoogleOAuth      config.GoogleOAuthConfig
	Slack            SlackConfig
	GrpcClient       config.GrpcClientConfig
	Postgres         config.PostgresConfig
	PostgresAsync    config.PostgresAsyncConfig
	Neo4j            config.Neo4jConfig
	Jaeger           tracing.JaegerConfig
	RabbitMQ         config.RabbitMQConfig
	OpenSRS          config.OpenSRSConfig
	Postmark         config.PostmarkConfig
	InternalServices InternalServices
}

type InternalServices struct {
	EnrichmentApi config.EnrichmentAPIConfig
	ValidationApi config.ValidationAPIConfig
}

// CustomerOSConfig holds Customer OS API configuration
type CustomerOSConfig struct {
	API    string `env:"CUSTOMER_OS_API,required"`
	APIKey string `env:"CUSTOMER_OS_API_KEY,required"`
}

// ServiceConfig holds main service configuration
type ServiceConfig struct {
	Port               string `env:"PORT,required"`
	PublicPath         string `env:"USER_ADMIN_API_PUBLIC_PATH,required"`
	CorsURL            string `env:"CORS_URL,required"`
	ProviderTenantName string `env:"PROVIDER_TENANT_NAME,required"`
	ProviderUsername   string `env:"PROVIDER_USERNAME,required"`
}

// SlackConfig holds Slack integration configuration
type SlackConfig struct {
	ClientID                      string `env:"SLACK_CLIENT_ID,required"`
	ClientSecret                  string `env:"SLACK_CLIENT_SECRET,required"`
	NotifyNewTenantRegisteredHook string `env:"SLACK_NOTIFY_NEW_TENANT_REGISTERED_WEBHOOK,required"`
}
