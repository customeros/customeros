package config

import (
	"flag"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstroredb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "ES microservice config path")
}

type rawConfig struct {
	ServiceName      string                         `mapstructure:"serviceName"`
	Logger           *logger.Config                 `mapstructure:"logger"`
	GRPC             GRPC                           `mapstructure:"grpc"`
	EventStoreConfig eventstroredb.EventStoreConfig `mapstructure:"eventStoreConfig"`
	Subscriptions    Subscriptions                  `mapstructure:"subscriptions"`
	Neo4j            Neo4j                          `mapstructure:"neo4j"`
	Jaeger           *tracing.Config                `mapstructure:"jaeger"`
}

// Validate the configuration file
func validate(cfg rawConfig) error {
	v := validator.New()

	if err := v.Struct(cfg); err != nil {
		return fmt.Errorf("invalid configuration file: %w", err)
	}

	// Perform additional validation here, if needed
	return nil
}

type GRPC struct {
	Port        string `mapstructure:"port"`
	Development bool   `mapstructure:"development"`
}

type Subscriptions struct {
	PoolSize                           int    `mapstructure:"poolSize" validate:"required,gte=0"`
	PhoneNumberPrefix                  string `mapstructure:"phoneNumberPrefix" validate:"required,gte=0"`
	EmailPrefix                        string `mapstructure:"emailPrefix" validate:"required,gte=0"`
	UserPrefix                         string `mapstructure:"userPrefix" validate:"required,gte=0"`
	ContactPrefix                      string `mapstructure:"contactPrefix" validate:"required,gte=0"`
	OrganizationPrefix                 string `mapstructure:"organizationPrefix" validate:"required,gte=0"`
	GraphProjectionGroupName           string `mapstructure:"graphProjectionGroupName" validate:"required,gte=0"`
	EmailValidationProjectionGroupName string `mapstructure:"emailValidationProjectionGroupName" validate:"required,gte=0"`
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

type Services struct {
	ValidationApi    string `mapstructure:"validationApi"`
	ValidationApiKey string `mapstructure:"validationApiKey"`
}

type Config struct {
	ServiceName      string
	Logger           *logger.Config
	GRPC             GRPC
	EventStoreConfig eventstroredb.EventStoreConfig
	Subscriptions    Subscriptions
	Neo4j            Neo4j
	Jaeger           *tracing.Config
	Services         Services
}

func InitConfig() (*Config, error) {
	// Load values from a .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

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

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err = viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	rawCfg := &rawConfig{}
	if err = viper.Unmarshal(&rawCfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}
	if err = validate(*rawCfg); err != nil {
		return nil, err
	}

	cfg := &Config{
		ServiceName:      rawCfg.ServiceName,
		Logger:           rawCfg.Logger,
		GRPC:             rawCfg.GRPC,
		EventStoreConfig: rawCfg.EventStoreConfig,
		Subscriptions:    rawCfg.Subscriptions,
		Neo4j:            rawCfg.Neo4j,
		Jaeger:           rawCfg.Jaeger,
	}

	if err := OverrideConfigWithEnvVars(cfg); err != nil {
		return nil, errors.Wrap(err, "OverrideConfigWithEnvVars")
	}

	return cfg, nil
}

// OverrideConfigWithEnvVars overrides the Config with environment variables
func OverrideConfigWithEnvVars(cfg *Config) error {
	if v, ok := os.LookupEnv(constants.EnvGrpcPort); ok {
		cfg.GRPC.Port = v
	}
	if v, ok := os.LookupEnv(constants.EnvValidationApiUrl); ok {
		cfg.Services.ValidationApi = v
	}
	if v, ok := os.LookupEnv(constants.EnvValidationApiKey); ok {
		cfg.Services.ValidationApiKey = v
	}
	if v, ok := os.LookupEnv(constants.EnvEventStoreConnectionString); ok {
		cfg.EventStoreConfig.ConnectionString = v
	}
	if v, ok := os.LookupEnv(constants.EnvJaegerHostPort); ok {
		cfg.Jaeger.HostPort = v
	}
	return nil
}
