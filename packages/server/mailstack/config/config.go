package config

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type AppConfig struct {
	APIPort           string `env:"PORT,required" envDefault:"12222"`
	APIKey            string `env:"API_KEY,required"`
	RabbitMQURL       string `env:"RABBITMQ_URL"`
	TrackingPublicUrl string `env:"TRACKING_PUBLIC_URL" envDefault:"https://custosmetrics.com"`
	Logger            *logger.Config
	Tracing           *tracing.JaegerConfig
}

type MailstackDatabaseConfig struct {
	Host            string `env:"MAILSTACK_POSTGRES_HOST,required"`
	Port            string `env:"MAILSTACK_POSTGRES_PORT,required"`
	User            string `env:"MAILSTACK_POSTGRES_USER,required"`
	DBName          string `env:"MAILSTACK_POSTGRES_DB_NAME,required"`
	Password        string `env:"MAILSTACK_POSTGRES_PASSWORD,required"`
	MaxConn         int    `env:"MAILSTACK_POSTGRES_DB_MAX_CONN"`
	MaxIdleConn     int    `env:"MAILSTACK_POSTGRES_DB_MAX_IDLE_CONN"`
	ConnMaxLifetime int    `env:"MAILSTACK_POSTGRES_DB_CONN_MAX_LIFETIME"`
	LogLevel        string `env:"MAILSTACK_POSTGRES_LOG_LEVEL" envDefault:"WARN"`
	SSLMode         string `env:"MAILSTACK_POSTGRES_SSL_MODE"`
}

type OpenlineDatabaseConfig struct {
	Host            string `env:"OPENLINE_POSTGRES_HOST,required"`
	Port            string `env:"OPENLINE_POSTGRES_PORT,required"`
	User            string `env:"OPENLINE_POSTGRES_USER,required"`
	DBName          string `env:"OPENLINE_POSTGRES_DB_NAME,required"`
	Password        string `env:"OPENLINE_POSTGRES_PASSWORD,required"`
	MaxConn         int    `env:"OPENLINE_POSTGRES_DB_MAX_CONN"`
	MaxIdleConn     int    `env:"OPENLINE_POSTGRES_DB_MAX_IDLE_CONN"`
	ConnMaxLifetime int    `env:"OPENLINE_POSTGRES_DB_CONN_MAX_LIFETIME"`
	LogLevel        string `env:"OPENLINE_POSTGRES_LOG_LEVEL" envDefault:"WARN"`
	SSLMode         string `env:"OPENLINE_POSTGRES_SSL_MODE"`
}

type R2StorageConfig struct {
	AccountID             string `env:"CLOUDFLARE_R2_ACCOUNT_ID,required"`
	AccessKeyID           string `env:"CLOUDFLARE_R2_ACCESS_KEY_ID,required"`
	AccessKeySecret       string `env:"CLOUDFLARE_R2_ACCESS_KEY_SECRET,required"`
	EmailAttachmentBucket string `env:"BUCKET_NAME_EMAIL_ATTACHMENT" envDefault:"attachments"`
}
