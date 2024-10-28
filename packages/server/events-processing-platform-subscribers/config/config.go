package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstoredb"
)

type Config struct {
	ServiceName        string `env:"SERVICE_NAME" envDefault:"events-processing-platform"`
	Logger             logger.Config
	EventStoreConfig   eventstoredb.EventStoreConfig
	GrpcClientConfig   config.GrpcClientConfig
	Neo4j              config.Neo4jConfig
	Postgres           config.PostgresConfig
	Jaeger             tracing.JaegerConfig
	RabbitMQConfig     config.RabbitMQConfig
	Subscriptions      Subscriptions
	Services           Services
	EventNotifications EventNotifications
	Temporal           config.TemporalConfig
}

type Subscriptions struct {
	GraphSubscription                 GraphSubscription
	GraphLowPrioritySubscription      GraphLowPrioritySubscription
	PhoneNumberValidationSubscription PhoneNumberValidationSubscription
	LocationValidationSubscription    LocationValidationSubscription
	OrganizationSubscription          OrganizationSubscription
	OrganizationWebscrapeSubscription OrganizationWebscrapeSubscription
	ContractSubscription              ContractSubscription
	NotificationsSubscription         NotificationsSubscription
	InvoiceSubscription               InvoiceSubscription
	ReminderSubscription              ReminderSubscription
	EnrichSubscription                EnrichSubscription
	NotifyRealtimeSubscription        NotifyRealtimeSubscription
}

type GraphSubscription struct {
	Enabled              bool   `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_ENABLED" envDefault:"true"`
	GroupName            string `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_GROUP_NAME" envDefault:"graph-v4" validate:"required"`
	PoolSize             int    `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_POOL_SIZE" envDefault:"10" validate:"required,gte=0"`
	BufferSizeClient     uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
	CheckpointLowerBound int32  `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_CHECKPOINT_LOWER_BOUND" envDefault:"10" validate:"required,gte=0"`
}

type GraphLowPrioritySubscription struct {
	Enabled          bool     `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_LOW_PRIO_V2_ENABLED" envDefault:"true"`
	GroupName        string   `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_LOW_PRIO_V2_GROUP_NAME" envDefault:"graph-low-prio-v2" validate:"required"`
	PoolSize         int      `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_LOW_PRIO_V2_POOL_SIZE" envDefault:"5" validate:"required,gte=0"`
	BufferSizeClient uint32   `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_LOW_PRIO_V2_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
	Prefixes         []string `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_LOW_PRIO_V2_PREFIXES" envDefault:"organization-" validate:"required"`
}

type PhoneNumberValidationSubscription struct {
	Enabled          bool   `env:"EVENT_STORE_SUBSCRIPTIONS_PHONE_NUMBER_VALIDATION_ENABLED" envDefault:"true"`
	GroupName        string `env:"EVENT_STORE_SUBSCRIPTIONS_PHONE_NUMBER_VALIDATION_GROUP_NAME" envDefault:"phoneNumberValidation-v2" validate:"required"`
	Prefix           string `env:"EVENT_STORE_SUBSCRIPTIONS_PHONE_NUMBER_PREFIX" envDefault:"phone_number-" validate:"required"`
	PoolSize         int    `env:"EVENT_STORE_SUBSCRIPTIONS_PHONE_NUMBER_VALIDATION_POOL_SIZE" envDefault:"5" validate:"required,gte=0"`
	BufferSizeClient uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_PHONE_NUMBER_VALIDATION_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
}

type LocationValidationSubscription struct {
	Enabled          bool   `env:"EVENT_STORE_SUBSCRIPTIONS_LOCATION_VALIDATION_ENABLED" envDefault:"true"`
	GroupName        string `env:"EVENT_STORE_SUBSCRIPTIONS_LOCATION_VALIDATION_GROUP_NAME" envDefault:"locationValidation-v2" validate:"required"`
	Prefix           string `env:"EVENT_STORE_SUBSCRIPTIONS_LOCATION_PREFIX" envDefault:"location-" validate:"required"`
	PoolSize         int    `env:"EVENT_STORE_SUBSCRIPTIONS_LOCATION_VALIDATION_POOL_SIZE" envDefault:"5" validate:"required,gte=0"`
	BufferSizeClient uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_LOCATION_VALIDATION_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
}

type OrganizationSubscription struct {
	Enabled                      bool   `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_ENABLED" envDefault:"true"`
	GroupName                    string `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_GROUP_NAME" envDefault:"organization-v2" validate:"required"`
	Prefix                       string `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_PREFIX" envDefault:"organization-" validate:"required"`
	PoolSize                     int    `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_POOL_SIZE" envDefault:"5" validate:"required,gte=0"`
	BufferSizeClient             uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
	MessageTimeoutSec            int32  `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_MESSAGE_TIMEOUT" envDefault:"180" validate:"required,gte=0"`
	CheckpointLowerBound         int32  `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_CHECKPOINT_LOWER_BOUND" envDefault:"10" validate:"required,gte=0"`
	DeletePersistentSubscription bool   `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_DELETE_SUBSCRIPTION" envDefault:"false"`
}

// TODO replace with EnrichSubscription
type OrganizationWebscrapeSubscription struct {
	Enabled                      bool   `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_WEBSCRAPE_ENABLED" envDefault:"true"`
	GroupName                    string `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_WEBSCRAPE_GROUP_NAME" envDefault:"organizationWebscrape-v2" validate:"required"`
	Prefix                       string `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_WEBSCRAPE_PREFIX" envDefault:"organization-" validate:"required"`
	PoolSize                     int    `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_WEBSCRAPE_POOL_SIZE" envDefault:"5" validate:"required,gte=0"`
	BufferSizeClient             uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_WEBSCRAPE_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
	MessageTimeoutSec            int32  `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_WEBSCRAPE_MESSAGE_TIMEOUT" envDefault:"300" validate:"required,gte=0"`
	CheckpointLowerBound         int32  `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_WEBSCRAPE_CHECKPOINT_LOWER_BOUND" envDefault:"4" validate:"required,gte=0"`
	DeletePersistentSubscription bool   `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_WEBSCRAPE_DELETE_SUBSCRIPTION" envDefault:"false"`
}

type EnrichSubscription struct {
	Enabled                      bool   `env:"EVENT_STORE_SUBSCRIPTIONS_ENRICH_ENABLED" envDefault:"true"`
	GroupName                    string `env:"EVENT_STORE_SUBSCRIPTIONS_ENRICH_GROUP_NAME" envDefault:"enrich-v2" validate:"required"`
	PoolSize                     int    `env:"EVENT_STORE_SUBSCRIPTIONS_ENRICH_POOL_SIZE" envDefault:"5" validate:"required,gte=0"`
	BufferSizeClient             uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_ENRICH_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
	MessageTimeoutSec            int32  `env:"EVENT_STORE_SUBSCRIPTIONS_ENRICH_MESSAGE_TIMEOUT" envDefault:"300" validate:"required,gte=0"`
	CheckpointLowerBound         int32  `env:"EVENT_STORE_SUBSCRIPTIONS_ENRICH_CHECKPOINT_LOWER_BOUND" envDefault:"4" validate:"required,gte=0"`
	DeletePersistentSubscription bool   `env:"EVENT_STORE_SUBSCRIPTIONS_ENRICH_DELETE_SUBSCRIPTION" envDefault:"false"`
}

type ContractSubscription struct {
	Enabled           bool   `env:"EVENT_STORE_SUBSCRIPTIONS_CONTRACT_ENABLED" envDefault:"true"`
	GroupName         string `env:"EVENT_STORE_SUBSCRIPTIONS_CONTRACT_GROUP_NAME" envDefault:"contract-v2" validate:"required"`
	Prefix            string `env:"EVENT_STORE_SUBSCRIPTIONS_CONTRACT_PREFIX" envDefault:"contract-" validate:"required"`
	PoolSize          int    `env:"EVENT_STORE_SUBSCRIPTIONS_CONTRACT_POOL_SIZE" envDefault:"5" validate:"required,gte=0"`
	BufferSizeClient  uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_CONTRACT_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
	MessageTimeoutSec int32  `env:"EVENT_STORE_SUBSCRIPTIONS_CONTRACT_MESSAGE_TIMEOUT" envDefault:"120" validate:"required,gte=0"`
}

type NotificationsSubscription struct {
	Enabled          bool   `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFICATIONS_ENABLED" envDefault:"true"`
	GroupName        string `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFICATIONS_GROUP_NAME" envDefault:"notifications-v2.1" validate:"required"`
	PoolSize         int    `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFICATIONS_POOL_SIZE" envDefault:"5" validate:"required,gte=0"`
	BufferSizeClient uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFICATIONS_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
	IgnoreEvents     bool   `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFICATIONS_IGNORE_EVENTS" envDefault:"true"`
	RedirectUrl      string `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFICATIONS_REDIRECT_URL" envDefault:"https://app.openline.dev"`
}

type InvoiceSubscription struct {
	Enabled           bool   `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_ENABLED" envDefault:"true"`
	GroupName         string `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_GROUP_NAME" envDefault:"invoice-v2" validate:"required"`
	PoolSize          int    `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_POOL_SIZE" envDefault:"5" validate:"required,gte=0"`
	MessageTimeoutSec int32  `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_MESSAGE_TIMEOUT" envDefault:"300" validate:"required,gte=0"`
	BufferSizeClient  uint32 `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_CLIENT_BUFFER_SIZE" envDefault:"5" validate:"required,gte=0"`
	IgnoreEvents      bool   `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_IGNORE_EVENTS" envDefault:"false"`
	PdfConverterUrl   string `env:"EVENT_STORE_INVOICE_NOTIFICATIONS_PDF_CONVERTER_URL" envDefault:"http://localhost:11006"`
}

type ReminderSubscription struct {
	Enabled          bool   `env:"EVENT_STORE_SUBSCRIPTIONS_REMINDER_ENABLED" envDefault:"true"`
	GroupName        string `env:"EVENT_STORE_SUBSCRIPTIONS_REMINDER_GROUP_NAME" envDefault:"reminder-v2" validate:"required"`
	PoolSize         int    `env:"EVENT_STORE_SUBSCRIPTIONS_REMINDER_POOL_SIZE" envDefault:"5" validate:"required,gte=0"`
	BufferSizeClient uint32 `env:"EVENT_STORE_SUBSCRIPTIONS_REMINDER_CLIENT_BUFFER_SIZE" envDefault:"10" validate:"required,gte=0"`
	IgnoreEvents     bool   `env:"EVENT_STORE_SUBSCRIPTIONS_REMINDER_IGNORE_EVENTS" envDefault:"false"`
}

type NotifyRealtimeSubscription struct {
	GroupName string `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFY_REALTIME_GROUP_NAME" envDefault:"notifyRealtime-v2" validate:"required"`
	Prefix    string `env:"EVENT_STORE_SUBSCRIPTIONS_NOTIFY_REALTIME_PREFIX" envDefault:"event_completion-" validate:"required"`
}

type Services struct {
	EnrichmentApi struct {
		Url    string `env:"ENRICHMENT_API_URL" validate:"required"`
		ApiKey string `env:"ENRICHMENT_API_KEY" validate:"required"`
	}
	ValidationApi struct {
		Url    string `env:"VALIDATION_API_URL" validate:"required"`
		ApiKey string `env:"VALIDATION_API_KEY" validate:"required"`
	}
	Ai struct {
		ApiPath string `env:"AI_API_PATH,required" envDefault:"N/A"`
		ApiKey  string `env:"AI_API_KEY,required" envDefault:"N/A"`
	}
	OpenAi struct {
		Organization string `env:"OPENAI_ORGANIZATION,required" envDefault:""`
	}
	Anthropic struct {
		IndustryLookupPrompt1    string `env:"ANTHROPIC_INDUSTRY_LOOKUP_PROMPT,required" envDefault:"With next Global Industry Classification Standard (GICS) valid values: (Aerospace & Defense,Air Freight & Logistics,Automobile Components,Automobiles,Banks,Beverages,Biotechnology,Broadline Retail,Building Products,Capital Markets,Chemicals,Commercial Services & Supplies,Communications Equipment,Construction & Engineering,Construction Materials,Consumer Finance,Consumer Staples Distribution & Retail,Containers & Packaging,Diversified Consumer Services,Diversified REITs,Diversified Telecommunication Services,Distributors,Electric Utilities,Electrical Equipment,Electronic Equipment,Instruments & Components,Energy Equipment & Services,Entertainment,Financial Services,Food Products,Gas Utilities,Ground Transportation,Health Care Equipment & Supplies,Health Care Providers & Services,Health Care REITs,Health Care Technology,Hotel & Resort REITs,Hotels,Restaurants & Leisure,Household Durables,Household Products,Independent Power and Renewable Electricity Producers,Industrial Conglomerates,Industrial REITs,Insurance,Interactive Media & Services,Internet Software & Services,IT Services,Leisure Products,Life Sciences Tools & Services,Machinery,Marine Transportation,Media,Metals & Mining,Mortgage Real Estate Investment Trusts (REITs),Multi-Utilities,Office REITs,Oil,Gas & Consumable Fuels,Paper & Forest Products,Passenger Airlines,Personal Products,Pharmaceuticals,Professional Services,Real Estate Management & Development,Residential REITs,Retail REITs,Semiconductors & Semiconductor Equipment,Software,Specialized REITs,Specialty Retail,Technology Hardware,Storage & Peripherals,Textiles,Apparel & Luxury Goods,Tobacco,Trading Companies & Distributors,Transportation Infrastructure,Water Utilities,Wireless Telecommunication Services), provide appropriate industry mapping for (%s) and if do not see obvious mapping, provide appropriate GICS value from the input list based on other companies providing similar services. Finally if cannot map return just single word: Unknown"`
		IndustryLookupPrompt2    string `env:"ANTHROPIC_INDUSTRY_LOOKUP_PROMPT,required" envDefault:"What GICS value from following list (Aerospace & Defense,Air Freight & Logistics,Automobile Components,Automobiles,Banks,Beverages,Biotechnology,Broadline Retail,Building Products,Capital Markets,Chemicals,Commercial Services & Supplies,Communications Equipment,Construction & Engineering,Construction Materials,Consumer Finance,Consumer Staples Distribution & Retail,Containers & Packaging,Diversified Consumer Services,Diversified REITs,Diversified Telecommunication Services,Distributors,Electric Utilities,Electrical Equipment,Electronic Equipment,Instruments & Components,Energy Equipment & Services,Entertainment,Financial Services,Food Products,Gas Utilities,Ground Transportation,Health Care Equipment & Supplies,Health Care Providers & Services,Health Care REITs,Health Care Technology,Hotel & Resort REITs,Hotels,Restaurants & Leisure,Household Durables,Household Products,Independent Power and Renewable Electricity Producers,Industrial Conglomerates,Industrial REITs,Insurance,Interactive Media & Services,Internet Software & Services,IT Services,Leisure Products,Life Sciences Tools & Services,Machinery,Marine Transportation,Media,Metals & Mining,Mortgage Real Estate Investment Trusts (REITs),Multi-Utilities,Office REITs,Oil,Gas & Consumable Fuels,Paper & Forest Products,Passenger Airlines,Personal Products,Pharmaceuticals,Professional Services,Real Estate Management & Development,Residential REITs,Retail REITs,Semiconductors & Semiconductor Equipment,Software,Specialized REITs,Specialty Retail,Technology Hardware,Storage & Peripherals,Textiles,Apparel & Luxury Goods,Tobacco,Trading Companies & Distributors,Transportation Infrastructure,Water Utilities,Wireless Telecommunication Services) is chosen in next statement. Strictly provide the value only: %s"`
		EmailSummaryPrompt       string `env:"ANTHROPIC_EMAIL_SUMMARY_PROMPT,required" envDefault:"Make a 120 characters summary for this html email: %v"`
		EmailActionsItemsPrompt  string `env:"ANTHROPIC_EMAIL_ACTIONS_ITEMS_PROMPT,required" envDefault:"Give me the action points to be taken for the email. The criticality for the action points should be at least medium severity. return response in jSON format, key - \"items\", value - array of strings. The email is: %v"`
		LocationEnrichmentPrompt string `env:"ANTHROPIC_LOCATION_ENRICHMENT_PROMPT,required" envDefault:"Given the address '%s', please provide a JSON representation of the Location object with all available information. Use the following structure, filling in as many fields as possible based on the given address. If a field cannot be determined, omit it from the JSON output. Strictly return only the JSON.\n\n{\n    \"country\": \"string\",\n    \"countryCodeA2\": \"string\",\n    \"countryCodeA3\": \"string\",\n    \"region\": \"string\",\n    \"locality\": \"string\",\n    \"address\": \"string\",\n    \"address2\": \"string\",\n    \"zip\": \"string\",\n    \"addressType\": \"string\",\n    \"houseNumber\": \"string\",\n    \"postalCode\": \"string\",\n    \"plusFour\": \"string\",\n    \"commercial\": boolean,\n    \"predirection\": \"string\",\n    \"district\": \"string\",\n    \"street\": \"string\",\n    \"latitude\": number,\n    \"longitude\": number,\n    \"timeZone\": \"string\",\n    \"utcOffset\": number\n}"`
	}
	Novu struct {
		ApiKey string `env:"NOVU_API_KEY,required" envDefault:"N/A"`
	}
	FileStoreApiConfig fsc.FileStoreApiConfig
	CustomerOsApi      struct {
		ApiUrl string `env:"CUSTOMER_OS_API_URL" envDefault:"https://api.customeros.ai" validate:"required"`
	}
}

type EventNotifications struct {
	SlackConfig struct {
		InternalAlertsRegisteredWebhook string `env:"SLACK_INTERNAL_ALERTS_REGISTERED_WEBHOOK" envDefault:""`
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
