package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	commonconf "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/joho/godotenv"

	cronconf "github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/cron/config"
)

type CommonConfig struct {
	Enrow             commonconf.EnrowConfig
	Logger            logger.Config
	Jaeger            tracing.JaegerConfig
	OpenTelemetry     telemetry.OpenTelemetryConfig
	RabbitMQConfig    commonconf.RabbitMQConfig
	ScrubbyIo         commonconf.ScrubbyIoConfig
	Anthropic         commonconf.AnthropicConfig
	Groq              commonconf.GroqConfig
	Gemini            commonconf.GeminiConfig
	CustomerOsApi     commonconf.CustomerOsApiConfig
	BetterContact     commonconf.BetterContactConfig
	Postgres          commonconf.PostgresConfig
	PostgresAsync     commonconf.PostgresAsyncConfig
	Neo4j             commonconf.Neo4jConfig
	Mailsherpa        commonconf.MailSherpaApiConfig
	MailstackApi      commonconf.MailstackApiConfig
	FileStore         commonconf.FileStoreConfig
	PdfConverter      commonconf.PdfConverterConfig
	SlackConfig       commonconf.SlackConfig
	GoogleOAuthConfig commonconf.GoogleOAuthConfig
	AzureOAuthConfig  commonconf.AzureOAuthConfig
	NovuConfig        commonconf.NovuConfig
	JinaConfig        commonconf.JinaConfig
	AwsConfig         commonconf.AwsConfig
	QuickbooksConfig  commonconf.QuickbooksConfig
	OpensearchConfig  commonconf.OpensearchConfig
	TemporalConfig    commonconf.TemporalConfig
	ScrapinConfig     commonconf.ScrapinConfig
}

type AppConfig struct {
	Cron               cronconf.Config
	ProcessConfig      ProcessConfig
	EventNotifications EventNotifications
	Limits             Limits
}

type Config struct {
	Common *commonconf.CommonConfig // Common shared configurations
	App    AppConfig                // Application-specific configurations
}

type ProcessConfig struct {
	CycleInvoicingEnabled                  bool `env:"CYCLE_INVOICING_ENABLED" envDefault:"true"`
	OffCycleInvoicingEnabled               bool `env:"OFF_CYCLE_INVOICING_ENABLED" envDefault:"false"`
	DelayAutoPayInvoiceInMinutes           int  `env:"DELAY_AUTO_PAY_INVOICE_IN_MINUTES" envDefault:"5"`
	DelayRequestPaymentLinkInMinutes       int  `env:"DELAY_REQUEST_PAYMENT_LINK_IN_MINUTES" envDefault:"15"`
	RequestPaymentLinkLookBackWindowInDays int  `env:"REQUEST_PAYMENT_LINK_LOOK_BACK_WINDOW_IN_DAYS" envDefault:"5"`
	DelayGenerateOffCycleInvoiceInMinutes  int  `env:"DELAY_GENERATE_OFF_CYCLE_INVOICE_IN_MINUTES" envDefault:"60"`
}

type Limits struct {
	EmailsValidationLimit       int `env:"EMAILS_VALIDATION_LIMIT" envDefault:"25" required:"true"`
	BulkEmailsValidationThreads int `env:"BULK_EMAILS_VALIDATION_THREADS" envDefault:"6" required:"true"`
}

type EventNotifications struct {
	IntegrationAppEventWebhookUrls struct {
		GeneratePaymentLinkUrl string `env:"INVOICE_GENERATE_PAYMENT_LINK_URL" envDefault:"" required:"true"`
		InvoiceFinalizedUrl    string `env:"INVOICE_READY_URL" envDefault:"" required:"true"`
	}
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found. Proceeding with system environment variables. Error: %v", err)
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error loading app configuration: %+v", err)
	}
	cmnCfg := CommonConfig{}
	if err := env.Parse(&cmnCfg); err != nil {
		log.Fatalf("Error loading app configuration: %+v", err)
	}

	cfg.Common = &commonconf.CommonConfig{
		Infrastructure: commonconf.InfrastructureConfig{
			LoggerConfig:        cmnCfg.Logger,
			JaegerConfig:        cmnCfg.Jaeger,
			OpenTelemetryConfig: cmnCfg.OpenTelemetry,
			RabbitMQConfig:      cmnCfg.RabbitMQConfig,
			PostgresConfig:      cmnCfg.Postgres,
			PostgresAsyncConfig: cmnCfg.PostgresAsync,
			Neo4jConfig:         cmnCfg.Neo4j,
			GoogleOAuthConfig:   cmnCfg.GoogleOAuthConfig,
			AzureOAuthConfig:    cmnCfg.AzureOAuthConfig,
			OpensearchConfig:    cmnCfg.OpensearchConfig,
		},
		External: commonconf.ExternalServicesConfig{
			EnrowConfig:         cmnCfg.Enrow,
			ScrubbyIoConfig:     cmnCfg.ScrubbyIo,
			AnthropicConfig:     cmnCfg.Anthropic,
			GroqConfig:          cmnCfg.Groq,
			GeminiConfig:        cmnCfg.Gemini,
			BetterContactConfig: cmnCfg.BetterContact,
			SlackConfig:         cmnCfg.SlackConfig,
			NovuConfig:          cmnCfg.NovuConfig,
			JinaConfig:          cmnCfg.JinaConfig,
			QuickbooksConfig:    cmnCfg.QuickbooksConfig,
			TemporalConfig:      cmnCfg.TemporalConfig,
			ScrapinConfig:       cmnCfg.ScrapinConfig,
		},
		Internal: commonconf.InternalServicesConfig{
			CustomerOsApi:       cmnCfg.CustomerOsApi,
			MailSherpaApiConfig: cmnCfg.Mailsherpa,
			MailstackApiConfig:  cmnCfg.MailstackApi,
			FileStoreConfig:     cmnCfg.FileStore,
			PdfConverterConfig:  cmnCfg.PdfConverter,
		},
	}

	return &cfg
}
