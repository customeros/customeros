package config

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"

type Config struct {
	Neo4jDb struct {
		Target                string `env:"NEO4J_TARGET,required"`
		User                  string `env:"NEO4J_AUTH_USER,required,unset"`
		Pwd                   string `env:"NEO4J_AUTH_PWD,required,unset"`
		Realm                 string `env:"NEO4J_AUTH_REALM"`
		MaxConnectionPoolSize int    `env:"NEO4J_MAX_CONN_POOL_SIZE" envDefault:"100"`
		LogLevel              string `env:"NEO4J_LOG_LEVEL" envDefault:"WARNING"`
	}
	PostgresDb struct {
		Host            string `env:"POSTGRES_HOST,required"`
		Port            string `env:"POSTGRES_PORT,required"`
		User            string `env:"POSTGRES_USER,required,unset"`
		Db              string `env:"POSTGRES_DB,required"`
		Password        string `env:"POSTGRES_PASSWORD,required,unset"`
		MaxConn         int    `env:"POSTGRES_DB_MAX_CONN"`
		MaxIdleConn     int    `env:"POSTGRES_DB_MAX_IDLE_CONN"`
		ConnMaxLifetime int    `env:"POSTGRES_DB_CONN_MAX_LIFETIME"`
		LogLevel        string `env:"POSTGRES_LOG_LEVEL" envDefault:"WARN"`
	}

	GoogleOAuth struct {
		ClientId     string `env:"GOOGLE_OAUTH_CLIENT_ID,required"`
		ClientSecret string `env:"GOOGLE_OAUTH_CLIENT_SECRET,required"`
	}

	SyncData struct {
		CronSync  string `env:"CRON_SYNC" envDefault:"0 */1 * * * *"`
		BatchSize int64  `env:"BATCH_SIZE" envDefault:"100"`

		GoogleCalendarReadBatchSize int64  `env:"GOOGLE_CALENDAR_READ_BATCH_SIZE" envDefault:"100"`
		GoogleCalendarSyncStartDate string `env:"GOOGLE_CALENDAR_SYNC_START_DATE,required"`
		GoogleCalendarSyncStopDate  string `env:"GOOGLE_CALENDAR_SYNC_STOP_DATE,required"`
	}

	Logger logger.Config
}
