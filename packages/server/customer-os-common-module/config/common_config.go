package config

type CommonConfig struct {
	PostgresConfig      *PostgresConfig
	PostgresAsyncConfig *PostgresAsyncConfig
	Neo4jConfig         *Neo4jConfig
	GoogleOAuthConfig   *GoogleOAuthConfig
	AzureOAuthConfig    *AzureOAuthConfig
	GrpcClientConfig    *GrpcClientConfig
	RabbitMQConfig      *RabbitMQConfig

	InternalServices *InternalServices
	ExternalServices *ExternalServices
}

type InternalServices struct {
	EmailConfig      *EmailConfig
	FileStoreConfig  *FileStoreConfig
	MailstackConfig  *MailstackConfig
	CustomerOsApiUrl string `env:"CUSTOMER_OS_API_URL" envDefault:"https://api.customeros.ai"`
}

type ExternalServices struct {
	AnthropicConfig      *AnthropicConfig
	AnthropicPrompts     *AnthropicPrompts
	BetterContactConfig  *BetterContactConfig
	BrandfetchConfig     *BrandfetchConfig
	CloudflareConfig     *CloudflareConfig
	EnrowConfig          *EnrowConfig
	IntegrationAppConfig *IntegrationAppConfig
	IpDataConfig         *IpDataConfig
	IpHunterConfig       *IpHunterConfig
	NamecheapConfig      *NamecheapConfig
	NovuCofig            *NovuCofig
	OpenSRSConfig        *OpenSRSConfig
	PostmarkConfig       *PostmarkConfig
	ScrapinConfig        *ScrapinConfig
	ScrubbyIoConfig      *ScrubbyIoConfig
	SlackConfig          *SlackConfig
	SmartyConfig         *SmartyConfig
	SnitcherConfig       *SnitcherConfig
	StripeConfig         *StripeConfig
	TemporalConfig       *TemporalConfig
	TrueInboxConfig      *TrueInboxConfig
}
