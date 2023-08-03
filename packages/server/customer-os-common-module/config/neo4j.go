package config

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/config"
	"time"
)

func NewNeo4jDriver(cfg Neo4jConfig) (neo4j.DriverWithContext, error) {
	return neo4j.NewDriverWithContext(
		cfg.Target,
		neo4j.BasicAuth(cfg.User, cfg.Pwd, cfg.Realm),
		func(config *config.Config) {
			config.MaxConnectionPoolSize = cfg.MaxConnectionPoolSize
			config.ConnectionAcquisitionTimeout = time.Duration(cfg.ConnectionAcquisitionTimeoutSec) * time.Second
			config.Log = neo4j.ConsoleLogger(strToLogLevel(cfg.LogLevel))
		})
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
