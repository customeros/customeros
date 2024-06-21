package config

type GlobalConfig struct {
	PostgresConfig    *PostgresConfig
	Neo4jConfig       *Neo4jConfig
	GoogleOAuthConfig *GoogleOAuthConfig
	GrpcClientConfig  *GrpcClientConfig
	TemporalConfig    *TemporalConfig
}
