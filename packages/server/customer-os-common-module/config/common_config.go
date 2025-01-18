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
	GrpcClientConfig    GrpcClientConfig
	RabbitMQConfig      RabbitMQConfig
	JaegerConfig        tracing.JaegerConfig
	LoggerConfig        logger.Config
}

type InternalServicesConfig struct {
	EmailConfig         EmailConfig
	FileStoreConfig     FileStoreConfig
	MailstackConfig     MailstackConfig
	MailSherpaApiConfig MailSherpaApiConfig
	CustomerOsApi       CustomerOsApiConfig
}

type ExternalServicesConfig struct {
	AnthropicConfig      AnthropicConfig
	AnthropicPrompts     AnthropicPrompts
	BetterContactConfig  BetterContactConfig
	BrandfetchConfig     BrandfetchConfig
	CloudflareConfig     CloudflareConfig
	EnrowConfig          EnrowConfig
	IntegrationAppConfig IntegrationAppConfig
	IpDataConfig         IpDataConfig
	NamecheapConfig      NamecheapConfig
	NovuCofig            NovuCofig
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
}
