package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
)

type Config struct {
	Logger     logger.Config
	CustomerOS struct {
		CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
	}
	Service struct {
		ServerAddress string `env:"USER_ADMIN_API_SERVER_ADDRESS,required"`
		CorsUrl       string `env:"USER_ADMIN_API_CORS_URL,required"`
		ApiKey        string `env:"USER_ADMIN_API_KEY,required"`
	}
	EventsProcessingPlatform struct {
		EventsProcessingPlatformUrl    string `env:"EVENTS_PROCESSING_PLATFORM_URL" validate:"required"`
		EventsProcessingPlatformApiKey string `env:"EVENTS_PROCESSING_PLATFORM_API_KEY" validate:"required"`
	}
	GoogleOAuth struct {
		ClientId     string `env:"GOOGLE_OAUTH_CLIENT_ID,required"`
		ClientSecret string `env:"GOOGLE_OAUTH_CLIENT_SECRET,required"`
	}
	Slack struct {
		ClientId     string `env:"SLACK_CLIENT_ID,required"`
		ClientSecret string `env:"SLACK_CLIENT_SECRET,required"`
	}
	Postgres config.PostgresConfig
	Neo4j    config.Neo4jConfig
}
