package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/metrics"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"log"
)

type Config struct {
	ApiPort     string `env:"PORT" envDefault:"10000" validate:"required"`
	MetricsPort string `env:"PORT_METRICS" envDefault:"10000" validate:"required"`
	Logger      logger.Config
	GraphQL     struct {
		PlaygroundEnabled    bool `env:"GRAPHQL_PLAYGROUND_ENABLED" envDefault:"false"`
		FixedComplexityLimit int  `env:"GRAPHQL_FIXED_COMPLEXITY_LIMIT" envDefault:"200"`
	}
	Admin struct {
		Key string `env:"ADMIN_KEY,required"`
	}
	GrpcClientConfig config.GrpcClientConfig
	Postgres         config.PostgresConfig
	Neo4j            config.Neo4jConfig
	Jaeger           tracing.JaegerConfig
	Metrics          metrics.Config
	AppConfig        struct {
		TrackingPublicUrl      string `env:"TRACKING_PUBLIC_URL" envDefault:"https://custosmetrics.com"`
		InvoicePaidRedirectUrl string `env:"INVOICE_PAID_REDIRECT_URL" envDefault:"https://customeros.ai/payments/status/paid/"`
	}
	InternalServices struct {
		CustomerOsApiUrl   string `env:"CUSTOMER_OS_API_URL" envDefault:"https://api.customeros.ai" validate:"required"`
		ValidationApi      string `env:"VALIDATION_API" validate:"required"`
		ValidationApiKey   string `env:"VALIDATION_API_KEY" validate:"required"`
		EnrichmentApiUrl   string `env:"ENRICHMENT_API_URL" validate:"required"`
		EnrichmentApiKey   string `env:"ENRICHMENT_API_KEY" validate:"required"`
		FileStoreApiConfig fsc.FileStoreApiConfig
	}
	ExternalServices struct {
		IntegrationApp struct {
			WorkspaceKey                    string `env:"INTEGRATION_APP_WORKSPACE_KEY"`
			WorkspaceSecret                 string `env:"INTEGRATION_APP_WORKSPACE_SECRET"`
			ApiTriggerUrlCreatePaymentLinks string `env:"INTEGRATION_APP_API_TRIGGER_URL_CREATE_PAYMENT_LINKS"`
		}
		Namecheap struct {
			ApiKey      string `env:"NAMECHEAP_API_KEY" validate:"required"`
			ApiUser     string `env:"NAMECHEAP_API_USER" validate:"required"`
			ApiUsername string `env:"NAMECHEAP_API_USERNAME" validate:"required"`
		}
	}
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
