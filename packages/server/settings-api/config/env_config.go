package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
)

type Config struct {
	ApiPort              string `env:"PORT"`
	Logger               logger.Config
	Postgres             config.PostgresConfig
	Neo4j                config.Neo4jConfig
	EncodedEncryptionKey string `env:"ENCODED_ENCRYPTION_KEY"`
}
