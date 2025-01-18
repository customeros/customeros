package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/validator"
	"github.com/customeros/customeros/packages/server/events/eventstoredb"
)

type Config struct {
	ServiceName      string `env:"SERVICE_NAME" envDefault:"events-processing-platform"`
	Logger           logger.Config
	EventStoreConfig eventstoredb.EventStoreConfig
	GrpcClientConfig config.GrpcClientConfig
	Jaeger           tracing.JaegerConfig
	Subscriptions    Subscriptions

	CommonServices config.CommonConfig
}

type Subscriptions struct {
	GraphSubscription         GraphSubscription
	ContractSubscription      ContractSubscription
	NotificationsSubscription NotificationsSubscription
	InvoiceSubscription       InvoiceSubscription
}

type GraphSubscription struct {
	Enabled              bool   `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_ENABLED" envDefault:"true"`
	GroupName            string `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_GROUP_NAME" envDefault:"graph-v5" validate:"required"`
	PoolSize             int    `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_POOL_SIZE" envDefault:"10" validate:"required,gte=0"`
	BufferSizeClient     uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
	CheckpointLowerBound int32  `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_CHECKPOINT_LOWER_BOUND" envDefault:"10" validate:"required,gte=0"`
}

type ContractSubscription struct {
	Enabled           bool   `env:"EVENT_STORE_SUBSCRIPTIONS_CONTRACT_ENABLED" envDefault:"true"`
	GroupName         string `env:"EVENT_STORE_SUBSCRIPTIONS_CONTRACT_GROUP_NAME" envDefault:"contract-v3" validate:"required"`
	Prefix            string `env:"EVENT_STORE_SUBSCRIPTIONS_CONTRACT_PREFIX" envDefault:"contract-" validate:"required"`
	PoolSize          int    `env:"EVENT_STORE_SUBSCRIPTIONS_CONTRACT_POOL_SIZE" envDefault:"5" validate:"required,gte=0"`
	BufferSizeClient  uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_CONTRACT_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
	MessageTimeoutSec int32  `env:"EVENT_STORE_SUBSCRIPTIONS_CONTRACT_MESSAGE_TIMEOUT" envDefault:"120" validate:"required,gte=0"`
}

type NotificationsSubscription struct {
	Enabled          bool   `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFICATIONS_ENABLED" envDefault:"true"`
	GroupName        string `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFICATIONS_GROUP_NAME" envDefault:"notifications-v3" validate:"required"`
	PoolSize         int    `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFICATIONS_POOL_SIZE" envDefault:"5" validate:"required,gte=0"`
	BufferSizeClient uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFICATIONS_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
	IgnoreEvents     bool   `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFICATIONS_IGNORE_EVENTS" envDefault:"true"`
	RedirectUrl      string `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFICATIONS_REDIRECT_URL" envDefault:"https://app.openline.dev"`
}

type InvoiceSubscription struct {
	Enabled           bool   `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_ENABLED" envDefault:"true"`
	GroupName         string `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_GROUP_NAME" envDefault:"invoice-v3" validate:"required"`
	PoolSize          int    `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_POOL_SIZE" envDefault:"5" validate:"required,gte=0"`
	MessageTimeoutSec int32  `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_MESSAGE_TIMEOUT" envDefault:"300" validate:"required,gte=0"`
	BufferSizeClient  uint32 `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_CLIENT_BUFFER_SIZE" envDefault:"5" validate:"required,gte=0"`
	IgnoreEvents      bool   `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_IGNORE_EVENTS" envDefault:"false"`
	PdfConverterUrl   string `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_PDF_CONVERTER_URL" envDefault:"http://localhost:11006"`
}

func InitConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	if err = env.Parse(&cfg); err != nil {
		return nil, err
	}

	err = validator.GetValidator().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
