package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	ApiPort       string `env:"PORT"`
	ApiServiceUrl string `env:"SERVICE_URL"`
	MaxFileSizeMB int64  `env:"MAX_FILE_SIZE_MB"`

	AWS struct {
		Region string `env:"AWS_S3_REGION,required"`
		Bucket string `env:"AWS_S3_BUCKET,required"`
	}

	PostgresConfig      config.PostgresConfig
	PostgresAsyncConfig config.PostgresAsyncConfig

	Service struct {
		CustomerOsAPI                  string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey               string `env:"CUSTOMER_OS_API_KEY,required"`
		FileStoreAPIJwtSecret          string `env:"FILE_STORE_API_JWT_SECRET,required"`
		CloudflareImageUploadAccountId string `env:"CLOUDFLARE_IMAGE_UPLOAD_ACCOUNT_ID" envDefault:""`
		CloudflareImageUploadApiKey    string `env:"CLOUDFLARE_IMAGE_UPLOAD_API_KEY" envDefault:""`
		CloudflareImageUploadSignKey   string `env:"CLOUDFLARE_IMAGE_UPLOAD_SIGN_KEY" envDefault:""`
	}
	Logger logger.Config
	Neo4j  config.Neo4jConfig
	Jaeger tracing.JaegerConfig
}
