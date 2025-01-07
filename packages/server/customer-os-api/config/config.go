package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/metrics"
)

type Config struct {
	Server           ServerConfig
	GraphQL          GraphQLConfig
	Admin            AdminConfig
	Database         DatabaseConfig
	Messaging        MessagingConfig
	Observability    ObservabilityConfig
	App              AppConfig
	InternalServices InternalServicesConfig
	ExternalServices ExternalServicesConfig
	GrpcClient       config.GrpcClientConfig
}

type ServerConfig struct {
	ApiPort     string `env:"PORT" envDefault:"10000" validate:"required"`
	MetricsPort string `env:"PORT_METRICS" envDefault:"10000" validate:"required"`
	Logger      logger.Config
}

type GraphQLConfig struct {
	PlaygroundEnabled    bool `env:"GRAPHQL_PLAYGROUND_ENABLED" envDefault:"false"`
	FixedComplexityLimit int  `env:"GRAPHQL_FIXED_COMPLEXITY_LIMIT" envDefault:"200"`
}

type AdminConfig struct {
	Key string `env:"ADMIN_KEY,required"`
}

type DatabaseConfig struct {
	Postgres      config.PostgresConfig
	PostgresAsync config.PostgresAsyncConfig
	Neo4j         config.Neo4jConfig
}

type MessagingConfig struct {
	RabbitMQ config.RabbitMQConfig
}

type ObservabilityConfig struct {
	Jaeger  tracing.JaegerConfig
	Metrics metrics.Config
}

type CORSConfig struct {
	AllowOrigins []string `env:"ALLOW_ORIGINS" envDefault:"*"`
	AllowHeaders []string `env:"ALLOW_HEADERS" envDefault:"x-openline-username"`
}

type MailstackConfig struct {
	SupportedTlds []string `env:"MAILSTACK_SUPPORTED_TLDS" envDefault:"com"`
}

type InternalServicesConfig struct {
	CustomerOsApiUrl string `env:"CUSTOMER_OS_API_URL" envDefault:"https://api.customeros.ai" validate:"required"`
	ValidationApi    string `env:"VALIDATION_API" validate:"required"`
	ValidationApiKey string `env:"VALIDATION_API_KEY" validate:"required"`
	EnrichmentApiUrl string `env:"ENRICHMENT_API_URL" validate:"required"`
	EnrichmentApiKey string `env:"ENRICHMENT_API_KEY" validate:"required"`
	FileStoreApi     fsc.FileStoreApiConfig
}

type IntegrationAppConfig struct {
	WorkspaceKey                    string `env:"INTEGRATION_APP_WORKSPACE_KEY"`
	WorkspaceSecret                 string `env:"INTEGRATION_APP_WORKSPACE_SECRET"`
	ApiTriggerUrlCreatePaymentLinks string `env:"INTEGRATION_APP_API_TRIGGER_URL_CREATE_PAYMENT_LINKS"`
}

type ExternalServicesConfig struct {
	Anthropic      config.AnthropicConfig
	Cloudflare     config.CloudflareConfig
	IntegrationApp IntegrationAppConfig
	Namecheap      config.NamecheapConfig
	OpenSRS        config.OpenSRSConfig
	Postmark       config.PostmarkConfig
	Stripe         config.StripeConfig
}

type AppConfig struct {
	CORS                                   CORSConfig
	TrackingPublicUrl                      string `env:"TRACKING_PUBLIC_URL" envDefault:"https://custosmetrics.com"`
	InvoicePaidRedirectUrl                 string `env:"INVOICE_PAID_REDIRECT_URL" envDefault:"https://customeros.ai/payments/status/paid/"`
	Mailstack                              MailstackConfig
	DefaultGlobalOrgPrimaryDomainsInSearch []string `env:"DEFAULT_GLOBAL_ORG_PRIMARY_DOMAINS_IN_SEARCH" envDefault:"stripe.com,zapier.com,braintreepayments.com,discord.com,airtable.com,framer.com,gocardless.com,gong.io,intercom.com,linear.app,loom.com,mailchimp.com,monday.com,notion.so,brex.com,monzo.com,mercury.com,thebrowser.company,descript.com,ramp.com,pleo.io,scale.com,perplexity.ai,runwayml.com,togetherai.com,pulley.com,pitch.com,raycast.com,height.app,tailscale.com,elevenlabs.io,hume.ai,huggingface.co,rabbit.com,figma.com,superhuman.com,vercel.com"`
}

func InitConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Print("Error loading .env file")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
	}

	err := validator.GetValidator().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
