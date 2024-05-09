package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	ApiPort string `env:"PORT"`

	ReacherConfig ReacherConfig
	SmartyConfig  SmartyConfig
	HunterConfig  HunterConfig

	Postgres config.PostgresConfig
	Neo4j    config.Neo4jConfig
	Jaeger   tracing.JaegerConfig
	Logger   logger.Config
}

type ReacherConfig struct {
	ApiPath string `env:"REACHER_API_PATH,required"`
	Secret  string `env:"REACHER_SECRET,required"`
}

type SmartyConfig struct {
	AuthId    string `env:"SMARTY_AUTH_ID,required"`
	AuthToken string `env:"SMARTY_AUTH_TOKEN,required"`
}

type HunterConfig struct {
	ApiKey   string  `env:"HUNTER_IO_API_KEY,required"`
	ApiPath  string  `env:"HUNTER_IO_API_PATH,required"`
	MinScore float64 `env:"HUNTER_IO_MIN_SCORE,required,default:75.0"`
}
