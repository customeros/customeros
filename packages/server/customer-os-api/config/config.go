package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	commonconf "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/validator"
	"github.com/joho/godotenv"

	"github.com/customeros/customeros/packages/server/customer-os-api/metrics"
)

type Config struct {
	App    AppConfig
	Common *commonconf.CommonConfig
}

type CommonConfig struct {
	Logger         logger.Config
	Jaeger         tracing.JaegerConfig
	RabbitMQConfig commonconf.RabbitMQConfig
	Postgres       commonconf.PostgresConfig
	PostgresAsync  commonconf.PostgresAsyncConfig
	Neo4j          commonconf.Neo4jConfig
	Opensearch     commonconf.OpensearchConfig
	CosApi         commonconf.CustomerOsApiConfig

	Anthropic         commonconf.AnthropicConfig
	AzureOAuthConfig  commonconf.AzureOAuthConfig
	BetterContact     commonconf.BetterContactConfig
	Brandfetch        commonconf.BrandfetchConfig
	Cloudflare        commonconf.CloudflareConfig
	Deepseek          commonconf.DeepseekConfig
	FileStore         commonconf.FileStoreConfig
	GoogleOAuthConfig commonconf.GoogleOAuthConfig
	Groq              commonconf.GroqConfig
	Gemini            commonconf.GeminiConfig
	IntegrationApp    commonconf.IntegrationAppConfig
	IpData            commonconf.IpDataConfig
	Jina              commonconf.JinaConfig
	Mailsherpa        commonconf.MailSherpaApiConfig
	MailstackApi      commonconf.MailstackApiConfig
	Namecheap         commonconf.NamecheapConfig
	Novu              commonconf.NovuConfig
	OpenSRS           commonconf.OpenSRSConfig
	PdfConverter      commonconf.PdfConverterConfig
	Postmark          commonconf.PostmarkConfig
	QuickbooksConfig  commonconf.QuickbooksConfig
	Stripe            commonconf.StripeConfig
	Scrapin           commonconf.ScrapinConfig
	SlackConfig       commonconf.SlackConfig
	Temporal          commonconf.TemporalConfig
	CrustData         commonconf.CrustDataConfig
}

type AppConfig struct {
	AuthConfig                             AuthConfig
	CORS                                   CORSConfig
	EncodedEncryptionKey                   string   `env:"ENCODED_ENCRYPTION_KEY"`
	TrackingPublicUrl                      string   `env:"TRACKING_PUBLIC_URL" envDefault:"https://custosmetrics.com"`
	InvoicePaidRedirectUrl                 string   `env:"INVOICE_PAID_REDIRECT_URL" envDefault:"https://customeros.ai/payments/status/paid/"`
	DefaultGlobalOrgPrimaryDomainsInSearch []string `env:"DEFAULT_GLOBAL_ORG_PRIMARY_DOMAINS_IN_SEARCH" envDefault:"stripe.com,zapier.com,braintreepayments.com,discord.com,airtable.com,framer.com,gocardless.com,gong.io,intercom.com,linear.app,loom.com,mailchimp.com,monday.com,notion.so,brex.com,monzo.com,mercury.com,thebrowser.company,descript.com,ramp.com,pleo.io,scale.com,perplexity.ai,runwayml.com,togetherai.com,pulley.com,pitch.com,raycast.com,height.app,tailscale.com,elevenlabs.io,hume.ai,huggingface.co,rabbit.com,figma.com,superhuman.com,vercel.com"`
	ApiPort                                string   `env:"PORT" envDefault:"10000" validate:"required"`
	AppKey                                 string   `env:"APP_KEY" envDefault:"" validate:"required"`
	AppAdminKey                            string   `env:"APP_ADMIN_KEY" envDefault:"" validate:"required"`
	MetricsPort                            string   `env:"PORT_METRICS" envDefault:"10000" validate:"required"`
	GraphQL                                GraphQLConfig
	Observability                          ObservabilityConfig
}

// Config helpers
type AuthConfig struct {
	ProviderTenantName string `env:"PROVIDER_TENANT_NAME,required"`
	ProviderUsername   string `env:"PROVIDER_USERNAME,required"`
}

type GraphQLConfig struct {
	PlaygroundEnabled    bool `env:"GRAPHQL_PLAYGROUND_ENABLED" envDefault:"false"`
	FixedComplexityLimit int  `env:"GRAPHQL_FIXED_COMPLEXITY_LIMIT" envDefault:"200"`
}

type ObservabilityConfig struct {
	Jaeger  tracing.JaegerConfig
	Metrics metrics.Config
}

type CORSConfig struct {
	AllowOrigins []string `env:"ALLOW_ORIGINS" envDefault:"*"`
	AllowHeaders []string `env:"ALLOW_HEADERS" envDefault:"x-openline-username"`
}

func InitConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Print("Error loading .env file")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error loading app configuration: %+v", err)
	}
	err := validator.GetValidator().Struct(cfg.App)
	if err != nil {
		return nil, err
	}

	cmnCfg := CommonConfig{}
	if err := env.Parse(&cmnCfg); err != nil {
		log.Fatalf("Error loading app configuration: %+v", err)
	}

	cfg.Common = &commonconf.CommonConfig{
		Infrastructure: commonconf.InfrastructureConfig{
			AzureOAuthConfig:    cmnCfg.AzureOAuthConfig,
			GoogleOAuthConfig:   cmnCfg.GoogleOAuthConfig,
			JaegerConfig:        cmnCfg.Jaeger,
			LoggerConfig:        cmnCfg.Logger,
			Neo4jConfig:         cmnCfg.Neo4j,
			OpensearchConfig:    cmnCfg.Opensearch,
			PostgresConfig:      cmnCfg.Postgres,
			PostgresAsyncConfig: cmnCfg.PostgresAsync,
			RabbitMQConfig:      cmnCfg.RabbitMQConfig,
		},
		External: commonconf.ExternalServicesConfig{
			AnthropicConfig:      cmnCfg.Anthropic,
			BrandfetchConfig:     cmnCfg.Brandfetch,
			CloudflareConfig:     cmnCfg.Cloudflare,
			DeepseekConfig:       cmnCfg.Deepseek,
			NamecheapConfig:      cmnCfg.Namecheap,
			OpenSRSConfig:        cmnCfg.OpenSRS,
			JinaConfig:           cmnCfg.Jina,
			IpDataConfig:         cmnCfg.IpData,
			StripeConfig:         cmnCfg.Stripe,
			BetterContactConfig:  cmnCfg.BetterContact,
			PostmarkConfig:       cmnCfg.Postmark,
			ScrapinConfig:        cmnCfg.Scrapin,
			SlackConfig:          cmnCfg.SlackConfig,
			QuickbooksConfig:     cmnCfg.QuickbooksConfig,
			IntegrationAppConfig: cmnCfg.IntegrationApp,
			NovuConfig:           cmnCfg.Novu,
			TemporalConfig:       cmnCfg.Temporal,
			GroqConfig:           cmnCfg.Groq,
			GeminiConfig:         cmnCfg.Gemini,
			CrustDataConfig:      cmnCfg.CrustData,
		},
		Internal: commonconf.InternalServicesConfig{
			CustomerOsApi:       cmnCfg.CosApi,
			FileStoreConfig:     cmnCfg.FileStore,
			MailSherpaApiConfig: cmnCfg.Mailsherpa,
			MailstackApiConfig:  cmnCfg.MailstackApi,
			PdfConverterConfig:  cmnCfg.PdfConverter,
		},
	}

	return &cfg, nil
}
