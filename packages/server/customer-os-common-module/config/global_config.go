package config

type GlobalConfig struct {
	PostgresConfig    *PostgresConfig
	Neo4jConfig       *Neo4jConfig
	GoogleOAuthConfig *GoogleOAuthConfig
	AzureOAuthConfig  *AzureOAuthConfig
	GrpcClientConfig  *GrpcClientConfig
	TemporalConfig    *TemporalConfig
	RabbitMQConfig    *RabbitMQConfig
	NovuConfig        *NovuConfig

	// Customer OS
	InternalServices InternalServices
	ExternalServices ExternalServices
}

type InternalServices struct {
	UserAdminApiPublicPath string `env:"USER_ADMIN_API_PUBLIC_PATH,required" envDefault:"http://localhost:4001"`
	EnrichmentApiConfig    EnrichmentAPIConfig
	AiApiConfig            AiAPIConfig
}

type ExternalServices struct {
	OpenSRSConfig OpenSRSConfig
}
