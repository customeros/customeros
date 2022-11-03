package config

type Config struct {
	ApiPort string `env:"PORT"`
	GraphQL struct {
		PlaygroundEnabled bool `env:"GRAPHQL_PLAYGROUND_ENABLED" envDefault:"false"`
	}
	Neo4j struct {
		Target                string `env:"NEO4J_TARGET,required"`
		User                  string `env:"NEO4J_AUTH_USER,required,unset"`
		Pwd                   string `env:"NEO4J_AUTH_PWD,required,unset"`
		Realm                 string `env:"NEO4J_AUTH_REALM"`
		MaxConnectionPoolSize int    `env:"NEO4J_MAX_CONN_POOL_SIZE" envDefault:"100"`
	}
}
