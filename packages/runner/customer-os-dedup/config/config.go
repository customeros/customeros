package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	cronConfig "github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/cron/config"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"log"
)

type Config struct {
	Postgres      commonConfig.PostgresConfig
	Neo4j         commonConfig.Neo4jConfig
	Logger        logger.Config
	Jaeger        tracing.JaegerConfig
	Cron          cronConfig.Config
	Organizations struct {
		AtLeastPerTenant       int `env:"ORGANIZATIONS_DEDUP_AT_LEAST_PER_TENANT" envDefault:"3"`
		OrganizationsPerPrompt int `env:"ORGANIZATIONS_DEDUP_ORGANIZATIONS_PER_PROMPT" envDefault:"100"`
		ForceDedupEachDays     int `env:"ORGANIZATIONS_DEDUP_FORCE_DEDUP_EACH_DAYS" envDefault:"7"`
		Anthropic              struct {
			Enabled            bool   `env:"ORGANIZATIONS_DEDUP_ANTHROPIC_ENABLED" envDefault:"true"`
			PromptSuggestNames string `env:"ORGANIZATIONS_DEDUP_ANTHROPIC_PROMPT_SUGGEST_NAMES" envDefault:"I have a list of organizations, each with a unique identifier (UUID) and a name. Some of the organizations may have similar or duplicate names. Identify possible duplicated organizations with confidence score at least 0.7 out of 1 and provide response in json format (maximum 2 pairs) containing pairs of organization IDs (mandatory from input so can be mapped later), their names and confidence score.\nResponse sample: {\"duplicates\":[{\"first\":{\"id\":\"id from input\",\"name\":\"name from input\"},\"second\":{\"id\":\"id from input\",\"name\":\"name from input\"},\"confidence\":0.XX}]}\nInput list of organizations: %s"`
			PromptCompareOrgs  string `env:"ORGANIZATIONS_DEDUP_ANTHROPIC_PROMPT_COMPARE_ORGS" envDefault:"I have details of 2 organizations (companies). Analyze details of both organizations and confidence level between 0 and 1 if organizations are same and should be merged into a single one. Identify what organization is primary and what is secondary to be merged into primary. Provide response in json format containing primary organization id, secondary organization id that should be merged into primary and confidence score between 0 and 1. If confidence level is < 0.49 return empty json. Sample of the output: {\"primary\": \"input ID here\", \"secondary\": \"input ID here\", \"confidence\": 0.XX}. Input organization id: %v and details: %v. Another input organization id: %v and details: %v"`
		}
		OpenAI struct {
			Enabled            bool   `env:"ORGANIZATIONS_DEDUP_OPENAI_ENABLED" envDefault:"false"`
			PromptSuggestNames string `env:"ORGANIZATIONS_DEDUP_OPENAI_PROMPT_SUGGEST_NAMES" envDefault:"I have a list of organizations, each with a unique identifier (UUID) and a name. Some of the organizations may have similar or duplicate names. Identify possible duplicated organizations with confidence score at least 0.7 out of 1 and provide response in json format containing pairs of organization IDs (mandatory from input so can be mapped later), their names and confidence score.\nResponse sample: {\"duplicates\":[{\"first\":{\"id\":\"id from input\",\"name\":\"name from input\"},\"second\":{\"id\":\"id from input\",\"name\":\"name from input\"},\"confidence\":0.XX}]}\nInput list of organizations: %s"`
			PromptCompareOrgs  string `env:"ORGANIZATIONS_DEDUP_OPENAI_PROMPT_COMPARE_ORGS" envDefault:"I have details of 2 organizations (companies). Analyze details of both organizations and confidence level between 0 and 1 if organizations are same and should be merged into a single one. Identify what organization is primary and what is secondary to be merged into primary. Provide response in json format containing primary organization id, secondary organization id that should be merged into primary and confidence score between 0 and 1. If confidence level is < 0.49 return empty json. Sample of the output: {\"primary\": \"input ID here\", \"secondary\": \"input ID here\", \"confidence\": 0.XX}. Input organization id: %v and details: %v. Another input organization id: %v and details: %v"`
		}
	}
	Service struct {
		CustomerOsAdminAPI    string `env:"CUSTOMER_OS_ADMIN_API,required"`
		CustomerOsAdminAPIKey string `env:"CUSTOMER_OS_ADMIN_API_KEY,required"`
		Anthropic             struct {
			ApiPath string `env:"ANTHROPIC_API,required" envDefault:"N/A"`
			ApiKey  string `env:"ANTHROPIC_API_KEY,required" envDefault:"N/A"`
		}
		OpenAI struct {
			ApiPath string `env:"OPENAI_API,required" envDefault:"N/A"`
			ApiKey  string `env:"OPENAI_API_KEY,required" envDefault:"N/A"`
		}
	}
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
