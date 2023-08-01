package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstroredb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
)

type GRPC struct {
	Port        string `env:"GRPC_PORT" envDefault:":5001" validate:"required"`
	Development bool   `env:"GRPC_DEVELOPMENT" envDefault:"false"`
	ApiKey      string `env:"GRPC_API_KEY" validate:"required"`
}

type Subscriptions struct {
	GraphSubscription                 GraphSubscription
	EmailValidationSubscription       EmailValidationSubscription
	PhoneNumberValidationSubscription PhoneNumberValidationSubscription
	LocationValidationSubscription    LocationValidationSubscription
	OrganizationSubscription          OrganizationSubscription
}

type GraphSubscription struct {
	Enabled          bool   `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_ENABLED" envDefault:"false"`
	GroupName        string `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_GROUP_NAME" envDefault:"graph-v1" validate:"required"`
	PoolSize         int    `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_POOL_SIZE" envDefault:"4" validate:"required,gte=0"`
	BufferSizeClient uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
}

type EmailValidationSubscription struct {
	Enabled          bool   `env:"EVENT_STORE_SUBSCRIPTIONS_EMAIL_VALIDATION_ENABLED" envDefault:"true"`
	GroupName        string `env:"EVENT_STORE_SUBSCRIPTIONS_EMAIL_VALIDATION_GROUP_NAME" envDefault:"emailValidation-v1" validate:"required"`
	Prefix           string `env:"EVENT_STORE_SUBSCRIPTIONS_EMAIL_PREFIX" envDefault:"email-" validate:"required"`
	PoolSize         int    `env:"EVENT_STORE_SUBSCRIPTIONS_EMAIL_VALIDATION_POOL_SIZE" envDefault:"1" validate:"required,gte=0"`
	BufferSizeClient uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_EMAIL_VALIDATION_CLIENT_BUFFER_SIZE" envDefault:"20" validate:"required,gte=0"`
}

type PhoneNumberValidationSubscription struct {
	Enabled          bool   `env:"EVENT_STORE_SUBSCRIPTIONS_PHONE_NUMBER_VALIDATION_ENABLED" envDefault:"true"`
	GroupName        string `env:"EVENT_STORE_SUBSCRIPTIONS_PHONE_NUMBER_VALIDATION_GROUP_NAME" envDefault:"phoneNumberValidation-v1" validate:"required"`
	Prefix           string `env:"EVENT_STORE_SUBSCRIPTIONS_PHONE_NUMBER_PREFIX" envDefault:"phone_number-" validate:"required"`
	PoolSize         int    `env:"EVENT_STORE_SUBSCRIPTIONS_PHONE_NUMBER_VALIDATION_POOL_SIZE" envDefault:"1" validate:"required,gte=0"`
	BufferSizeClient uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_PHONE_NUMBER_VALIDATION_CLIENT_BUFFER_SIZE" envDefault:"20" validate:"required,gte=0"`
}

type LocationValidationSubscription struct {
	Enabled          bool   `env:"EVENT_STORE_SUBSCRIPTIONS_LOCATION_VALIDATION_ENABLED" envDefault:"true"`
	GroupName        string `env:"EVENT_STORE_SUBSCRIPTIONS_LOCATION_VALIDATION_GROUP_NAME" envDefault:"locationValidation-v1" validate:"required"`
	Prefix           string `env:"EVENT_STORE_SUBSCRIPTIONS_LOCATION_PREFIX" envDefault:"location-" validate:"required"`
	PoolSize         int    `env:"EVENT_STORE_SUBSCRIPTIONS_LOCATION_VALIDATION_POOL_SIZE" envDefault:"1" validate:"required,gte=0"`
	BufferSizeClient uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_LOCATION_VALIDATION_CLIENT_BUFFER_SIZE" envDefault:"20" validate:"required,gte=0"`
}

type OrganizationSubscription struct {
	Enabled                      bool   `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_ENABLED" envDefault:"true"`
	GroupName                    string `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_GROUP_NAME" envDefault:"organization-v1" validate:"required"`
	Prefix                       string `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_PREFIX" envDefault:"organization-" validate:"required"`
	PoolSize                     int    `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_POOL_SIZE" envDefault:"1" validate:"required,gte=0"`
	BufferSizeClient             uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_CLIENT_BUFFER_SIZE" envDefault:"2" validate:"required,gte=0"`
	MessageTimeoutSec            int32  `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_MESSAGE_TIMEOUT" envDefault:"180" validate:"required,gte=0"`
	CheckpointLowerBound         int32  `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_CHECKPOINT_LOWER_BOUND" envDefault:"4" validate:"required,gte=0"`
	DeletePersistentSubscription bool   `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_DELETE_SUBSCRIPTION" envDefault:"false"`
}

type Neo4j struct {
	Target                          string `env:"NEO4J_TARGET" validate:"required"`
	User                            string `env:"NEO4J_AUTH_USER,unset" validate:"required"`
	Pwd                             string `env:"NEO4J_AUTH_PWD,unset" validate:"required"`
	Realm                           string `env:"NEO4J_AUTH_REALM"`
	MaxConnectionPoolSize           int    `env:"NEO4J_MAX_CONN_POOL_SIZE" envDefault:"100"`
	ConnectionAcquisitionTimeoutSec int    `env:"NEO4J_CONN_ACQUISITION_TIMEOUT_SEC" envDefault:"60"`
	LogLevel                        string `env:"NEO4J_LOG_LEVEL" envDefault:"WARNING"`
}

type Services struct {
	ValidationApi     string `env:"VALIDATION_API" validate:"required"`
	ValidationApiKey  string `env:"VALIDATION_API_KEY" validate:"required"`
	ScrapingBeeApiKey string `env:"SCRAPING_BEE_API_KEY" validate:"required"`
	OpenAi            struct {
		ApiPath             string `env:"OPENAI_API_PATH,required" envDefault:"N/A"`
		ApiKey              string `env:"OPENAI_API_KEY,required" envDefault:"N/A"`
		ScrapeCompanyPrompt string `env:"SCRAPE_COMPANY_PROMPT,required" envDefault:"Analyze the following text from a company website.
                         {{text}}
                         Analyze the text and respond (in English) as defined below:
                         {
                           companyName:  the name of the company,
                           market: One of the following options: B2B, B2C, or Marketplace,
                           industry: Industry per the Global Industry Classification Standard (GISB),
                           industryGroup: Industry Group per the Global Industry Classification Standard (GISB),
                           subIndustry: Sub-industry per the Global Industry Classification Standard (GISB),
                           targetAudience: analysis of the company'1s target audience,
                           valueProposition: analysis of the company's core value proposition,
                         }"`
		ScrapeDataPrompt string `env:"SCRAPE_DATA_PROMPT,required" envDefault:"The following is data scraped from a website:  Please combine and format the data into a clean json response

                      {{ANALYSIS}}

                      website: {{DOMAIN_URL}}

                      {{SOCIALS}}

                      --------

                      Put the data above in the following JSON structure

                      {{JSON_STRUCTURE}}

                      If you do not have data for a given key, do not return it as part of the JSON object.

                      Ensure before you output that your response is valid JSON.  If it is not valid JSON, do your best to fix the formatting to align to valid JSON."`
	}
	Anthropic struct {
		ApiPath               string `env:"ANTHROPIC_API,required" envDefault:"N/A"`
		ApiKey                string `env:"ANTHROPIC_API_KEY,required" envDefault:"N/A"`
		IndustryLookupPrompt1 string `env:"ANTHROPIC_INDUSTRY_LOOKUP_PROMPT,required" envDefault:"With next Global Industry Classification Standard (GICS) valid values: (Aerospace & Defense,Air Freight & Logistics,Automobile Components,Automobiles,Banks,Beverages,Biotechnology,Broadline Retail,Building Products,Capital Markets,Chemicals,Commercial Services & Supplies,Communications Equipment,Construction & Engineering,Construction Materials,Consumer Finance,Consumer Staples Distribution & Retail,Containers & Packaging,Diversified Consumer Services,Diversified REITs,Diversified Telecommunication Services,Distributors,Electric Utilities,Electrical Equipment,Electronic Equipment,Instruments & Components,Energy Equipment & Services,Entertainment,Financial Services,Food Products,Gas Utilities,Ground Transportation,Health Care Equipment & Supplies,Health Care Providers & Services,Health Care REITs,Health Care Technology,Hotel & Resort REITs,Hotels,Restaurants & Leisure,Household Durables,Household Products,Independent Power and Renewable Electricity Producers,Industrial Conglomerates,Industrial REITs,Insurance,Interactive Media & Services,Internet Software & Services,IT Services,Leisure Products,Life Sciences Tools & Services,Machinery,Marine Transportation,Media,Metals & Mining,Mortgage Real Estate Investment Trusts (REITs),Multi-Utilities,Office REITs,Oil,Gas & Consumable Fuels,Paper & Forest Products,Passenger Airlines,Personal Products,Pharmaceuticals,Professional Services,Real Estate Management & Development,Residential REITs,Retail REITs,Semiconductors & Semiconductor Equipment,Software,Specialized REITs,Specialty Retail,Technology Hardware,Storage & Peripherals,Textiles,Apparel & Luxury Goods,Tobacco,Trading Companies & Distributors,Transportation Infrastructure,Water Utilities,Wireless Telecommunication Services), provide appropriate industry mapping for (%s) and if do not see obvious mapping, provide appropriate GICS value from the input list based on other companies providing similar services. Finally if cannot map return just single word: Unknown"`
		IndustryLookupPrompt2 string `env:"ANTHROPIC_INDUSTRY_LOOKUP_PROMPT,required" envDefault:"What GICS value from following list (Aerospace & Defense,Air Freight & Logistics,Automobile Components,Automobiles,Banks,Beverages,Biotechnology,Broadline Retail,Building Products,Capital Markets,Chemicals,Commercial Services & Supplies,Communications Equipment,Construction & Engineering,Construction Materials,Consumer Finance,Consumer Staples Distribution & Retail,Containers & Packaging,Diversified Consumer Services,Diversified REITs,Diversified Telecommunication Services,Distributors,Electric Utilities,Electrical Equipment,Electronic Equipment,Instruments & Components,Energy Equipment & Services,Entertainment,Financial Services,Food Products,Gas Utilities,Ground Transportation,Health Care Equipment & Supplies,Health Care Providers & Services,Health Care REITs,Health Care Technology,Hotel & Resort REITs,Hotels,Restaurants & Leisure,Household Durables,Household Products,Independent Power and Renewable Electricity Producers,Industrial Conglomerates,Industrial REITs,Insurance,Interactive Media & Services,Internet Software & Services,IT Services,Leisure Products,Life Sciences Tools & Services,Machinery,Marine Transportation,Media,Metals & Mining,Mortgage Real Estate Investment Trusts (REITs),Multi-Utilities,Office REITs,Oil,Gas & Consumable Fuels,Paper & Forest Products,Passenger Airlines,Personal Products,Pharmaceuticals,Professional Services,Real Estate Management & Development,Residential REITs,Retail REITs,Semiconductors & Semiconductor Equipment,Software,Specialized REITs,Specialty Retail,Technology Hardware,Storage & Peripherals,Textiles,Apparel & Luxury Goods,Tobacco,Trading Companies & Distributors,Transportation Infrastructure,Water Utilities,Wireless Telecommunication Services) is chosen in next statement. Strictly provide the value only: %s"`
	}
}

type Config struct {
	ServiceName      string `env:"SERVICE_NAME" envDefault:"events-processing-platform"`
	Logger           logger.Config
	GRPC             GRPC
	EventStoreConfig eventstroredb.EventStoreConfig
	Subscriptions    Subscriptions
	Neo4j            Neo4j
	Jaeger           tracing.Config
	Services         Services
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
