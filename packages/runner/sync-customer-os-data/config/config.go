package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	commconf "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"log"
)

type Config struct {
	Neo4jDb           commconf.Neo4jConfig
	Postgres          commconf.PostgresConfig
	AirbytePostgresDb struct {
		Host            string `env:"DB_AIRBYTE_HOST,required"`
		Port            int    `env:"DB_AIRBYTE_PORT,required"`
		Pwd             string `env:"DB_AIRBYTE_PWD,required,unset"`
		User            string `env:"DB_AIRBYTE_USER,required"`
		MaxConn         int    `env:"DB_AIRBYTE_MAX_CONN"`
		MaxIdleConn     int    `env:"DB_AIRBYTE_MAX_IDLE_CONN"`
		ConnMaxLifetime int    `env:"DB_AIRBYTE_CONN_MAX_LIFETIME"`
		Name            string `env:"DB_AIRBYTE_NAME,required"`
	}
	SyncCustomerOsData struct {
		TimeoutAfterTaskRun int `env:"TIMEOUT_AFTER_TASK_RUN_SEC" envDefault:"60"`
		BatchSize           int `env:"SYNC_CUSTOMER_OS_DATA_BATCH_SIZE" envDefault:"10"`
	}

	GrpcClientConfig commconf.GrpcClientConfig
	SyncToEventStore struct {
		BatchSize           int  `env:"SYNC_TO_EVENT_STORE_BATCH_SIZE" envDefault:"100"`
		Enabled             bool `env:"SYNC_TO_EVENT_STORE_ENABLED" envDefault:"false"`
		TimeoutAfterTaskRun int  `env:"SYNC_TO_EVENT_STORE_TIMEOUT_AFTER_TASK_RUN_SEC" envDefault:"30"`
		Emails              struct {
			Enabled   bool `env:"SYNC_TO_EVENT_STORE_EMAILS_ENABLED" envDefault:"true"`
			BatchSize int  `env:"SYNC_TO_EVENT_STORE_EMAILS_BATCH_SIZE" envDefault:"-1"`
		}
		PhoneNumbers struct {
			Enabled   bool `env:"SYNC_TO_EVENT_STORE_PHONE_NUMBERS_ENABLED" envDefault:"true"`
			BatchSize int  `env:"SYNC_TO_EVENT_STORE_PHONE_NUMBERS_BATCH_SIZE" envDefault:"-1"`
		}
		Locations struct {
			Enabled   bool `env:"SYNC_TO_EVENT_STORE_LOCATIONS_ENABLED" envDefault:"true"`
			BatchSize int  `env:"SYNC_TO_EVENT_STORE_LOCATIONS_BATCH_SIZE" envDefault:"-1"`
		}
		Contacts struct {
			Enabled   bool `env:"SYNC_TO_EVENT_STORE_CONTACTS_ENABLED" envDefault:"true"`
			BatchSize int  `env:"SYNC_TO_EVENT_STORE_CONTACTS_BATCH_SIZE" envDefault:"-1"`
		}
		Organizations struct {
			Enabled                      bool `env:"SYNC_TO_EVENT_STORE_ORGANIZATIONS_ENABLED" envDefault:"true"`
			BatchSize                    int  `env:"SYNC_TO_EVENT_STORE_ORGANIZATIONS_BATCH_SIZE" envDefault:"-1"`
			OrganizationDomainsBatchSize int  `env:"SYNC_TO_EVENT_STORE_ORGANIZATIONS_DOMAINS_BATCH_SIZE" envDefault:"20"`
		}
	}
	Logger  logger.Config
	Service struct {
		CustomerOsWebhooksAPI    string `env:"CUSTOMER_OS_WEBHOOKS_API,required"`
		CustomerOsWebhooksAPIKey string `env:"CUSTOMER_OS_WEBHOOKS_API_KEY,required"`
	}
	Jaeger tracing.Config
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("Failed loading .env file")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
	}

	return &cfg
}
