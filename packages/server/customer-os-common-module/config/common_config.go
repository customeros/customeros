package config

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type CommonConfig struct {
	Infrastructure InfrastructureConfig
	Internal       InternalServicesConfig
	External       ExternalServicesConfig
}

type InfrastructureConfig struct {
	PostgresConfig      PostgresConfig
	PostgresAsyncConfig PostgresAsyncConfig
	Neo4jConfig         Neo4jConfig
	GoogleOAuthConfig   GoogleOAuthConfig
	AzureOAuthConfig    AzureOAuthConfig
	RabbitMQConfig      RabbitMQConfig
	OpensearchConfig    OpensearchConfig
	JaegerConfig        tracing.JaegerConfig
	LoggerConfig        logger.Config
}

type InternalServicesConfig struct {
	EmailConfig         EmailConfig
	MailstackConfig     MailstackConfig
	MailSherpaApiConfig MailSherpaApiConfig
	CustomerOsApi       CustomerOsApiConfig
	FileStoreConfig     FileStoreConfig
	PdfConverterConfig  PdfConverterConfig
}

type ExternalServicesConfig struct {
	AnthropicConfig      AnthropicConfig
	BetterContactConfig  BetterContactConfig
	BrandfetchConfig     BrandfetchConfig
	CloudflareConfig     CloudflareConfig
	DeepseekConfig       DeepseekConfig
	EnrowConfig          EnrowConfig
	GroqConfig           GroqConfig
	GeminiConfig         GeminiConfig
	IntegrationAppConfig IntegrationAppConfig
	IpDataConfig         IpDataConfig
	JinaConfig           JinaConfig
	NamecheapConfig      NamecheapConfig
	NovuConfig           NovuConfig
	OpenSRSConfig        OpenSRSConfig
	PostmarkConfig       PostmarkConfig
	ScrapinConfig        ScrapinConfig
	ScrubbyIoConfig      ScrubbyIoConfig
	SlackConfig          SlackConfig
	SnitcherConfig       SnitcherConfig
	SmartyConfig         SmartyConfig
	StripeConfig         StripeConfig
	TemporalConfig       TemporalConfig
	TrueInboxConfig      TrueInboxConfig
	QuickbooksConfig     QuickbooksConfig
	CrustDataConfig      CrustDataConfig
}
