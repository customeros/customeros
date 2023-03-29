package config

import (
	"flag"
	"fmt"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-common/constants"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstroredb"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/logger"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "ES microservice config path")
}

type Config struct {
	ServiceName      string           `mapstructure:"serviceName"`
	Logger           *logger.Config   `mapstructure:"logger"`
	GRPC             GRPC             `mapstructure:"grpc"`
	MongoCollections MongoCollections `mapstructure:"mongoCollections"`
	/*Jaeger           *tracing.Config                `mapstructure:"jaeger"`*/
	EventStoreConfig eventstroredb.EventStoreConfig `mapstructure:"eventStoreConfig"`
	Subscriptions    Subscriptions                  `mapstructure:"subscriptions"`
	Http             Http                           `mapstructure:"http"`
}

type GRPC struct {
	Port        string `mapstructure:"port"`
	Development bool   `mapstructure:"development"`
}

type MongoCollections struct {
	Orders string `mapstructure:"orders" validate:"required"`
}

type Subscriptions struct {
	PoolSize                   int    `mapstructure:"poolSize" validate:"required,gte=0"`
	OrderPrefix                string `mapstructure:"orderPrefix" validate:"required,gte=0"`
	MongoProjectionGroupName   string `mapstructure:"mongoProjectionGroupName" validate:"required,gte=0"`
	ElasticProjectionGroupName string `mapstructure:"elasticProjectionGroupName" validate:"required,gte=0"`
}

type Http struct {
	Port                string   `mapstructure:"port" validate:"required"`
	Development         bool     `mapstructure:"development"`
	BasePath            string   `mapstructure:"basePath" validate:"required"`
	OrdersPath          string   `mapstructure:"ordersPath" validate:"required"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/config/config.yaml", getwd)
		}
	}

	cfg := &Config{}

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	grpcPort := os.Getenv(constants.GrpcPort)
	if grpcPort != "" {
		cfg.GRPC.Port = grpcPort
	}

	/*	jaegerAddr := os.Getenv(constants.JaegerHostPort)
		if jaegerAddr != "" {
			cfg.Jaeger.HostPort = jaegerAddr
		}
	*/
	eventStoreConnectionString := os.Getenv(constants.EventStoreConnectionString)
	if eventStoreConnectionString != "" {
		cfg.EventStoreConfig.ConnectionString = eventStoreConnectionString
	}

	return cfg, nil
}
