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
	ServiceName      string                         `mapstructure:"serviceName"`
	Logger           *logger.Config                 `mapstructure:"logger"`
	GRPC             GRPC                           `mapstructure:"grpc"`
	EventStoreConfig eventstroredb.EventStoreConfig `mapstructure:"eventStoreConfig"`
	Subscriptions    Subscriptions                  `mapstructure:"subscriptions"`
	Neo4j            Neo4j                          `mapstructure:"neo4j"`
	/*Jaeger           *tracing.Config                `mapstructure:"jaeger"`*/
}

type GRPC struct {
	Port        string `mapstructure:"port"`
	Development bool   `mapstructure:"development"`
}

type Subscriptions struct {
	PoolSize                        int    `mapstructure:"poolSize" validate:"required,gte=0"`
	PhoneNumberPrefix               string `mapstructure:"phoneNumberPrefix" validate:"required,gte=0"`
	GraphProjectionGroupName        string `mapstructure:"graphProjectionGroupName" validate:"required,gte=0"`
	DataEnricherProjectionGroupName string `mapstructure:"dataEnricherProjectionGroupName" validate:"required,gte=0"`
}

type Neo4j struct {
	Target                          string `mapstructure:"target" validate:"required"`
	User                            string `mapstructure:"user" validate:"required"`
	Pwd                             string `mapstructure:"password" validate:"required"` // FIXME alexb implement unset
	Realm                           string `mapstructure:"realm"`
	MaxConnectionPoolSize           int    `mapstructure:"maxConnectionPoolSize" validate:"required"`
	ConnectionAcquisitionTimeoutSec int    `mapstructure:"connectionAcquisitionTimeoutSec" validate:"required"`
	LogLevel                        string `mapstructure:"logLevel"`
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
