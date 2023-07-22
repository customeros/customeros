package config

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/logger"
)

func NewDriver(log logger.Logger, cfg *Config) (*neo4j.DriverWithContext, error) {
	log.Infof("Connecting to neo4j database %s", cfg.Neo4jDb.Target)
	neo4jDriver, err := neo4j.NewDriverWithContext(
		cfg.Neo4jDb.Target,
		neo4j.BasicAuth(cfg.Neo4jDb.User, cfg.Neo4jDb.Pwd, cfg.Neo4jDb.Realm),
		func(config *config.Config) {
			config.MaxConnectionPoolSize = cfg.Neo4jDb.MaxConnectionPoolSize
			config.Log = neo4j.ConsoleLogger(strToLogLevel(cfg.Neo4jDb.LogLevel))
		})
	return &neo4jDriver, err
}

func strToLogLevel(str string) neo4j.LogLevel {
	switch str {
	case "ERROR":
		return neo4j.ERROR
	case "INFO":
		return neo4j.INFO
	case "DEBUG":
		return neo4j.DEBUG
	}
	return neo4j.WARNING
}
