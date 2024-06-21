package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
)

type Config struct {
	Neo4jDb    config.Neo4jConfig
	PostgresDb config.PostgresConfig

	AuthConfig config.GoogleOAuthConfig

	SyncData struct {
		CronSync  string `env:"CRON_SYNC" envDefault:"0 */1 * * * *"`
		BatchSize int64  `env:"BATCH_SIZE" envDefault:"100"`

		GoogleCalendarReadBatchSize int64  `env:"GOOGLE_CALENDAR_READ_BATCH_SIZE" envDefault:"100"`
		GoogleCalendarSyncStartDate string `env:"GOOGLE_CALENDAR_SYNC_START_DATE,required"`
		GoogleCalendarSyncStopDate  string `env:"GOOGLE_CALENDAR_SYNC_STOP_DATE,required"`
	}

	Logger logger.Config
}
