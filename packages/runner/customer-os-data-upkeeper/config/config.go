package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	cronconf "github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/cron/config"
	commconf "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"log"
)

type CommonConfig struct {
	Enrow            commconf.EnrowConfig
	Logger           logger.Config
	Jaeger           tracing.JaegerConfig
	GrpcClientConfig commconf.GrpcClientConfig
	RabbitMQConfig   commconf.RabbitMQConfig
	ScrubbyIo        commconf.ScrubbyIoConfig
	Anthropic        commconf.AnthropicConfig
	CustomerOsApi    commconf.CustomerOsApiConfig
	BetterContact    commconf.BetterContactConfig
	Postgres         commconf.PostgresConfig
	PostgresAsync    commconf.PostgresAsyncConfig
	Neo4j            commconf.Neo4jConfig
}

type AppConfig struct {
	Cron               cronconf.Config
	ProcessConfig      ProcessConfig
	EventNotifications EventNotifications
	Limits             Limits
}

type Config struct {
	Common *commconf.CommonConfig
	App    AppConfig
}

type ProcessConfig struct {
	CycleInvoicingEnabled                    bool `env:"CYCLE_INVOICING_ENABLED" envDefault:"true"`
	OffCycleInvoicingEnabled                 bool `env:"OFF_CYCLE_INVOICING_ENABLED" envDefault:"false"`
	DelaySendPayInvoiceNotificationInMinutes int  `env:"DELAY_SEND_PAY_INVOICE_NOTIFICATION_IN_MINUTES" envDefault:"60"`
	RetrySendPayInvoiceNotificationDays      int  `env:"RETRY_SEND_PAY_INVOICE_NOTIFICATION_DAYS" envDefault:"5"`
	DelayAutoPayInvoiceInMinutes             int  `env:"DELAY_AUTO_PAY_INVOICE_IN_MINUTES" envDefault:"5"`
	DelayRequestPaymentLinkInMinutes         int  `env:"DELAY_REQUEST_PAYMENT_LINK_IN_MINUTES" envDefault:"15"`
	RequestPaymentLinkLookBackWindowInDays   int  `env:"REQUEST_PAYMENT_LINK_LOOK_BACK_WINDOW_IN_DAYS" envDefault:"5"`
	DelayGenerateCycleInvoiceInMinutes       int  `env:"DELAY_GENERATE_CYCLE_INVOICE_IN_MINUTES" envDefault:"240"`
	DelayGenerateOffCycleInvoiceInMinutes    int  `env:"DELAY_GENERATE_OFF_CYCLE_INVOICE_IN_MINUTES" envDefault:"60"`
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
		log.Print("Failed loading .env file")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
	}
	cmnCfg := CommonConfig{}
	if err := env.Parse(&cmnCfg); err != nil {
		log.Fatalf("%+v", err)
	}

	cfg.Common = &commconf.CommonConfig{
		Infrastructure: commconf.InfrastructureConfig{
			LoggerConfig:        cmnCfg.Logger,
			JaegerConfig:        cmnCfg.Jaeger,
			GrpcClientConfig:    cmnCfg.GrpcClientConfig,
			RabbitMQConfig:      cmnCfg.RabbitMQConfig,
			PostgresConfig:      cmnCfg.Postgres,
			PostgresAsyncConfig: cmnCfg.PostgresAsync,
			Neo4jConfig:         cmnCfg.Neo4j,
		},
		External: commconf.ExternalServicesConfig{
			EnrowConfig:         cmnCfg.Enrow,
			ScrubbyIoConfig:     cmnCfg.ScrubbyIo,
			AnthropicConfig:     cmnCfg.Anthropic,
			BetterContactConfig: cmnCfg.BetterContact,
		},
		Internal: commconf.InternalServicesConfig{
			CustomerOsApi: cmnCfg.CustomerOsApi,
		},
	}

	return &cfg
}
