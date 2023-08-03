package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	Service struct {
		CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
		FileStoreAPI     string `env:"FILE_STORE_API,required"`
		FileStoreAPIKey  string `env:"FILE_STORE_API_KEY,required"`
		ServerAddress    string `env:"COMMS_API_SERVER_ADDRESS,required"`
		CorsUrl          string `env:"COMMS_API_CORS_URL,required"`
	}
	Mail struct {
		ApiKey string `env:"COMMS_API_MAIL_API_KEY,required"`
	}
	WebChat struct {
		PingInterval int `env:"WEBSOCKET_PING_INTERVAL"`
	}
	VCon struct {
		ApiKey string `env:"COMMS_API_VCON_API_KEY,required"`
	}
	Postgres struct {
		Host            string `env:"POSTGRES_HOST,required"`
		Port            string `env:"POSTGRES_PORT,required"`
		User            string `env:"POSTGRES_USER,required,unset"`
		Db              string `env:"POSTGRES_DB,required"`
		Password        string `env:"POSTGRES_PASSWORD,required,unset"`
		MaxConn         int    `env:"POSTGRES_DB_MAX_CONN"`
		MaxIdleConn     int    `env:"POSTGRES_DB_MAX_IDLE_CONN"`
		ConnMaxLifetime int    `env:"POSTGRES_DB_CONN_MAX_LIFETIME"`
	}
	WebRTC struct {
		AuthSecret   string `env:"WEBRTC_AUTH_SECRET,required"`
		TTL          int    `env:"WEBRTC_AUTH_TTL,required"`
		PingInterval int    `env:"WEBSOCKET_PING_INTERVAL"`
	}
	CalCom struct {
		CalComWebhookSecret string `env:"CALCOM_SECRET,required"`
	}
	Redis struct {
		Host string `env:"REDIS_HOST,required"`
	}
	Jaeger tracing.Config
	Logger logger.Config
}
