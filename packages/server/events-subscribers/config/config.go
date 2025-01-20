package config

import (
	"github.com/caarlos0/env/v6"
	commonconf "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Common *commonconf.CommonConfig
	App    AppConfig
}

type AppConfig struct {
	// Add config here for events-subscribers that are not included in commonconf.CommonConfig
}

type CommonConfig struct {
	Logger        logger.Config
	Jaeger        tracing.JaegerConfig
	Postgres      commonconf.PostgresConfig
	PostgresAsync commonconf.PostgresAsyncConfig
	Neo4j         commonconf.Neo4jConfig
	RabbitMQ      commonconf.RabbitMQConfig
	GrpcClient    commonconf.GrpcClientConfig
	MailSherpaApi commonconf.MailSherpaApiConfig
	BetterContact commonconf.BetterContactConfig
	Scrapin       commonconf.ScrapinConfig
	Snitcher      commonconf.SnitcherConfig
	Anthropic     commonconf.AnthropicConfig
	Novu          commonconf.NovuCofig
	Namecheap     commonconf.NamecheapConfig
	OpenSrs       commonconf.OpenSRSConfig
	Cloudflare    commonconf.CloudflareConfig
	TrueInbox     commonconf.TrueInboxConfig
	Brandfetch    commonconf.BrandfetchConfig
	Enrow         commonconf.EnrowConfig
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

	cfg.Common = &commonconf.CommonConfig{
		Infrastructure: commonconf.InfrastructureConfig{
			LoggerConfig:        cmnCfg.Logger,
			JaegerConfig:        cmnCfg.Jaeger,
			GrpcClientConfig:    cmnCfg.GrpcClient,
			PostgresConfig:      cmnCfg.Postgres,
			PostgresAsyncConfig: cmnCfg.PostgresAsync,
			Neo4jConfig:         cmnCfg.Neo4j,
			RabbitMQConfig:      cmnCfg.RabbitMQ,
		},
		Internal: commonconf.InternalServicesConfig{
			MailSherpaApiConfig: cmnCfg.MailSherpaApi,
		},
		External: commonconf.ExternalServicesConfig{
			BetterContactConfig: cmnCfg.BetterContact,
			ScrapinConfig:       cmnCfg.Scrapin,
			SnitcherConfig:      cmnCfg.Snitcher,
			AnthropicConfig:     cmnCfg.Anthropic,
			NovuCofig:           cmnCfg.Novu,
			NamecheapConfig:     cmnCfg.Namecheap,
			OpenSRSConfig:       cmnCfg.OpenSrs,
			CloudflareConfig:    cmnCfg.Cloudflare,
			TrueInboxConfig:     cmnCfg.TrueInbox,
			BrandfetchConfig:    cmnCfg.Brandfetch,
			EnrowConfig:         cmnCfg.Enrow,
		},
	}

	return &cfg
}
