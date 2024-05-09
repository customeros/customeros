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

	Ai struct {
		ApiPath string `env:"AI_API_PATH,required" envDefault:"WARN"`
		ApiKey  string `env:"AI_API_KEY,required" envDefault:"WARN"`
	}

	Anthropic struct {
		SummaryPrompt    string `env:"ANTHROPIC_SUMMARY_PROMPT,required" envDefault:"WARN"`
		ActionItemsPromp string `env:"ANTHROPIC_ACTION_ITEMS_PROMPT,required" envDefault:"WARN"`
	}

	SyncData struct {
		CronSync string `env:"CRON_SYNC" envDefault:"0 */1 * * * *"`
	}

	Jaeger tracing.JaegerConfig
	Logger logger.Config
}
