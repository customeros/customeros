package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	Neo4jDb    config.Neo4jConfig
	PostgresDb config.PostgresConfig

	GrpcClientConfig config.GrpcClientConfig

	Anthropic struct {
		ApiPath          string `env:"ANTHROPIC_API_PATH,required" envDefault:"WARN"`
		ApiKey           string `env:"ANTHROPIC_API_KEY,required" envDefault:"WARN"`
		SummaryPrompt    string `env:"ANTHROPIC_SUMMARY_PROMPT,required" envDefault:"WARN"`
		ActionItemsPromp string `env:"ANTHROPIC_ACTION_ITEMS_PROMPT,required" envDefault:"WARN"`
	}

	OpenAi struct {
		ApiPath string `env:"OPENAI_API_PATH,required" envDefault:"WARN"`
		ApiKey  string `env:"OPENAI_API_KEY,required" envDefault:"WARN"`
	}

	SyncData struct {
		CronSync string `env:"CRON_SYNC" envDefault:"0 */1 * * * *"`
	}

	Jaeger tracing.JaegerConfig
	Logger logger.Config
}
