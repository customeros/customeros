package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	commonconf "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/joho/godotenv"
)

type Config struct {
	Common commonconf.CommonConfig
	App    AppConfig
}

type AppConfig struct {
	// Add config here for events-subscribers that are not included in commonconf.CommonConfig
}

type CommonConfig struct {
	Logger           logger.Config
	Jaeger           tracing.JaegerConfig
	Postgres         commonconf.PostgresConfig
	PostgresAsync    commonconf.PostgresAsyncConfig
	Neo4j            commonconf.Neo4jConfig
	RabbitMQ         commonconf.RabbitMQConfig
	OpensearchConfig commonconf.OpensearchConfig
	MailSherpaApi    commonconf.MailSherpaApiConfig
	MaistackApi      commonconf.MailstackApiConfig
	BetterContact    commonconf.BetterContactConfig
	Scrapin          commonconf.ScrapinConfig
	Snitcher         commonconf.SnitcherConfig
	Anthropic        commonconf.AnthropicConfig
	Novu             commonconf.NovuConfig
	Namecheap        commonconf.NamecheapConfig
	OpenSrs          commonconf.OpenSRSConfig
	Cloudflare       commonconf.CloudflareConfig
	TrueInbox        commonconf.TrueInboxConfig
	Brandfetch       commonconf.BrandfetchConfig
	Enrow            commonconf.EnrowConfig
	QuickbooksConfig commonconf.QuickbooksConfig
	PdfConverter     commonconf.PdfConverterConfig
	FileStore        commonconf.FileStoreConfig
	SlackConfig      commonconf.SlackConfig
	CustomerOsApi    commonconf.CustomerOsApiConfig
	IntegrationApp   commonconf.IntegrationAppConfig
	Temporal         commonconf.TemporalConfig
	JinaConfig       commonconf.JinaConfig
	GeminiConfig     commonconf.GeminiConfig
	Groq             commonconf.GroqConfig
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

	cfg.Common = commonconf.CommonConfig{
		Infrastructure: commonconf.InfrastructureConfig{
			LoggerConfig:        cmnCfg.Logger,
			JaegerConfig:        cmnCfg.Jaeger,
			PostgresConfig:      cmnCfg.Postgres,
			PostgresAsyncConfig: cmnCfg.PostgresAsync,
			Neo4jConfig:         cmnCfg.Neo4j,
			RabbitMQConfig:      cmnCfg.RabbitMQ,
			OpensearchConfig:    cmnCfg.OpensearchConfig,
		},
		Internal: commonconf.InternalServicesConfig{
			MailSherpaApiConfig: cmnCfg.MailSherpaApi,
			MailstackApiConfig:  cmnCfg.MaistackApi,
			PdfConverterConfig:  cmnCfg.PdfConverter,
			FileStoreConfig:     cmnCfg.FileStore,
			CustomerOsApi:       cmnCfg.CustomerOsApi,
		},
		External: commonconf.ExternalServicesConfig{
			BetterContactConfig:  cmnCfg.BetterContact,
			ScrapinConfig:        cmnCfg.Scrapin,
			SnitcherConfig:       cmnCfg.Snitcher,
			AnthropicConfig:      cmnCfg.Anthropic,
			NovuConfig:           cmnCfg.Novu,
			NamecheapConfig:      cmnCfg.Namecheap,
			OpenSRSConfig:        cmnCfg.OpenSrs,
			CloudflareConfig:     cmnCfg.Cloudflare,
			TrueInboxConfig:      cmnCfg.TrueInbox,
			BrandfetchConfig:     cmnCfg.Brandfetch,
			EnrowConfig:          cmnCfg.Enrow,
			QuickbooksConfig:     cmnCfg.QuickbooksConfig,
			SlackConfig:          cmnCfg.SlackConfig,
			IntegrationAppConfig: cmnCfg.IntegrationApp,
			TemporalConfig:       cmnCfg.Temporal,
			JinaConfig:           cmnCfg.JinaConfig,
			GeminiConfig:         cmnCfg.GeminiConfig,
			GroqConfig:           cmnCfg.Groq,
		},
	}

	return &cfg
}
