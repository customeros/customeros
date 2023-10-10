package config

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func NewDriver(cfg *Config) (neo4j.DriverWithContext, error) {
	return neo4j.NewDriverWithContext(
		cfg.Neo4j.Target,
		neo4j.BasicAuth(cfg.Neo4j.User, cfg.Neo4j.Pwd, cfg.Neo4j.Realm),
		func(config *neo4j.Config) {
			config.MaxConnectionPoolSize = cfg.Neo4j.MaxConnectionPoolSize
			config.Log = neo4j.ConsoleLogger(strToLogLevel(cfg.Neo4j.LogLevel))
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
