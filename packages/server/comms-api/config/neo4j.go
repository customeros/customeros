package config

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func NewDriver(cfg *Config) (neo4j.DriverWithContext, error) {
	return neo4j.NewDriverWithContext(
		cfg.Neo4jConfig.Target,
		neo4j.BasicAuth(cfg.Neo4jConfig.User, cfg.Neo4jConfig.Pwd, cfg.Neo4jConfig.Realm),
		func(config *neo4j.Config) {
			config.MaxConnectionPoolSize = cfg.Neo4jConfig.MaxConnectionPoolSize
			config.Log = neo4j.ConsoleLogger(strToLogLevel(cfg.Neo4jConfig.LogLevel))
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
