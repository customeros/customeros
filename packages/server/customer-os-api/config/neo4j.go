package config

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func NewDriver(cfg *Config) (neo4j.Driver, error) {
	return neo4j.NewDriver(
		cfg.Neo4j.Target,
		neo4j.BasicAuth(cfg.Neo4j.User, cfg.Neo4j.Pwd, cfg.Neo4j.Realm),
		func(config *neo4j.Config) {
			config.MaxConnectionPoolSize = cfg.Neo4j.MaxConnectionPoolSize
		})
}
