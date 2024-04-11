package config

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"

type Config struct {
	ApiPort string `env:"PORT"`

	ReacherApiPath string `env:"REACHER_API_PATH,required"`
	ReacherSecret  string `env:"REACHER_SECRET,required"`

	Postgres config.PostgresConfig
	Neo4j    config.Neo4jConfig

	Smarty struct {
		AuthId    string `env:"SMARTY_AUTH_ID,required"`
		AuthToken string `env:"SMARTY_AUTH_TOKEN,required"`
	}
}
