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
	RabbitMQConfig   config.RabbitMQConfig
	Jaeger           tracing.JaegerConfig
	Metrics          metrics.Config
	AppConfig        struct {
		AllowOrigins           []string `env:"ALLOW_ORIGINS" envDefault:"*"`
		AllowHeaders           []string `env:"ALLOW_HEADERS" envDefault:"x-openline-username"`
		TrackingPublicUrl      string   `env:"TRACKING_PUBLIC_URL" envDefault:"https://custosmetrics.com"`
		InvoicePaidRedirectUrl string   `env:"INVOICE_PAID_REDIRECT_URL" envDefault:"https://customeros.ai/payments/status/paid/"`
		Mailstack              struct {
			SupportedTlds []string `env:"MAILSTACK_SUPPORTED_TLDS" envDefault:"com"`
		}
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
			Url                   string  `env:"NAMECHEAP_URL" envDefault:"https://api.namecheap.com/xml.response" validate:"required"`
			ApiKey                string  `env:"NAMECHEAP_API_KEY" validate:"required"`
			ApiUser               string  `env:"NAMECHEAP_API_USER" validate:"required"`
			ApiUsername           string  `env:"NAMECHEAP_API_USERNAME" validate:"required"`
			ApiClientIp           string  `env:"NAMECHEAP_API_CLIENT_IP" validate:"required"`
			MaxPrice              float64 `env:"NAMECHEAP_MAX_PRICE" envDefault:"20.0" validate:"required"`
			Years                 int     `env:"NAMECHEAP_YEARS" envDefault:"1" validate:"required"`
			RegistrantFirstName   string  `env:"NAMECHEAP_REGISTRANT_FIRST_NAME" validate:"required"`
			RegistrantLastName    string  `env:"NAMECHEAP_REGISTRANT_LAST_NAME" validate:"required"`
			RegistrantCompanyName string  `env:"NAMECHEAP_REGISTRANT_COMPANY_NAME" validate:"required"`
			RegistrantJobTitle    string  `env:"NAMECHEAP_REGISTRANT_JOB_TITLE" validate:"required"`
			RegistrantAddress1    string  `env:"NAMECHEAP_REGISTRANT_ADDRESS1" validate:"required"`
			RegistrantCity        string  `env:"NAMECHEAP_REGISTRANT_CITY" validate:"required"`
			RegistrantState       string  `env:"NAMECHEAP_REGISTRANT_STATE" validate:"required"`
			RegistrantZIP         string  `env:"NAMECHEAP_REGISTRANT_ZIP" validate:"required"`
			RegistrantCountry     string  `env:"NAMECHEAP_REGISTRANT_COUNTRY" validate:"required"`
			RegistrantPhoneNumber string  `env:"NAMECHEAP_REGISTRANT_PHONE_NUMBER" validate:"required"`
			RegistrantEmail       string  `env:"NAMECHEAP_REGISTRANT_EMAIL" validate:"required"`
		}
		Cloudflare struct {
			Url    string `env:"CLOUDFLARE_URL" envDefault:"https://api.cloudflare.com/client/v4" validate:"required"`
			ApiKey string `env:"CLOUDFLARE_API_KEY" validate:"required"`
			Email  string `env:"CLOUDFLARE_API_EMAIL" validate:"required"`
		}
		OpenSRSConfig config.OpenSRSConfig
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
