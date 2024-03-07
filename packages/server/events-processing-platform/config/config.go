package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstroredb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
)

type Config struct {
	ServiceName        string `env:"SERVICE_NAME" envDefault:"events-processing-platform"`
	Logger             logger.Config
	EventStoreConfig   eventstroredb.EventStoreConfig
	Neo4j              config.Neo4jConfig
	Postgres           config.PostgresConfig
	Jaeger             tracing.JaegerConfig
	GRPC               GRPC
	Services           Services
	Utils              Utils
	EventNotifications EventNotifications
}

type GRPC struct {
	Port        string `env:"GRPC_PORT" envDefault:":5001" validate:"required"`
	Development bool   `env:"GRPC_DEVELOPMENT" envDefault:"false"`
	ApiKey      string `env:"GRPC_API_KEY" validate:"required"`
}

type Services struct {
	ValidationApi                  string `env:"VALIDATION_API" validate:"required"`
	ValidationApiKey               string `env:"VALIDATION_API_KEY" validate:"required"`
	EventsProcessingPlatformUrl    string `env:"EVENTS_PROCESSING_PLATFORM_URL" validate:"required"`
	EventsProcessingPlatformApiKey string `env:"EVENTS_PROCESSING_PLATFORM_API_KEY" validate:"required"`
	ScrapingBeeApiKey              string `env:"SCRAPING_BEE_API_KEY" validate:"required"`
	CoreSignalApiKey               string `env:"CORE_SIGNAL_API_KEY" validate:"required" envDefault:"{}"`
	PromptJsonSchema               string `env:"PROMPT_JSON_SCHEMA" validate:"required" envDefault:"{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type": "object",
		"properties": {
		  "companyName": {
			"type": "string",
			"description": "the name of the company"
		  },
		  "market": {
			"type": "string",
			"description": "One of the following options: [B2B, B2C, or Marketplace]"
		  },
		  "industry": {
			"type": "string",
			"description": "Industry category per the Global Industry Classification Standard (GISB)"
		  },
		  "industryGroup": {
			"type": "string",
			"description": "Industry Group per the Global Industry Classification Standard (GISB)"
		  },
		  "subIndustry": {
			"type": "string",
			"description": "Sub-industry category per the Global Industry Classification Standard (GISB)"
		  },
		  "targetAudience": {
			"type": "string",
			"description": "analysis of the company's target audience"
		  },
		  "valueProposition": {
			"type": "string",
			"description": "analysis of the company's core value proposition"
		  }
		},
		"required": ["companyName", "market", "valueProposition", "industry"],
		"additionalProperties": false
	  }"`
	OpenAi struct {
		ApiPath             string `env:"OPENAI_API_PATH,required" envDefault:"N/A"`
		ApiKey              string `env:"OPENAI_API_KEY,required" envDefault:"N/A"`
		Organization        string `env:"OPENAI_ORGANIZATION,required" envDefault:""`
		ScrapeCompanyPrompt string `env:"SCRAPE_COMPANY_PROMPT,required" envDefault:"Analyze the text below and return the complete schema {{jsonschema}}\n\nTEXT\n{{text}}"`
		ScrapeDataPrompt    string `env:"SCRAPE_DATA_PROMPT,required" envDefault:"The following is data scraped from a website:  Please combine and format the data into a clean json response

                      {{ANALYSIS}}

                      website: {{DOMAIN_URL}}

                      {{SOCIALS}}

                      --------

                      Put the data above in the following JSON structure

                      {{JSON_STRUCTURE}}

                      If you do not have data for a given key, do not return it as part of the JSON object.

                      Ensure before you output that your response is valid JSON.  If it is not valid JSON, do your best to fix the formatting to align to valid JSON."`
	}
	Novu struct {
		ApiKey string `env:"NOVU_API_KEY,required" envDefault:"N/A"`
	}
	FileStoreApiConfig fsc.FileStoreApiConfig
}

type Utils struct {
	RetriesOnOptimisticLockException int `env:"UTILS_RETRIES_ON_OPTIMISTIC_LOCK" envDefault:"5"`
}

type EventNotifications struct {
	EndPoints struct {
		InvoiceReady string `env:"INVOICE_READY_URL" envDefault:""`
	}
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
