package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	Logger   logger.Config
	Postgres config.PostgresConfig
	Neo4j    config.Neo4jConfig
	Jaeger   tracing.JaegerConfig
	RabbitMQ config.RabbitMQConfig
	Services Services
}

type Services struct {
	EnrichmentApi struct {
		Url    string `env:"ENRICHMENT_API_URL" validate:"required"`
		ApiKey string `env:"ENRICHMENT_API_KEY" validate:"required"`
	}
	//ValidationApi struct {
	//	Url    string `env:"VALIDATION_API_URL" validate:"required"`
	//	ApiKey string `env:"VALIDATION_API_KEY" validate:"required"`
	//}
	//CustomerOsApi struct {
	//	ApiUrl string `env:"CUSTOMER_OS_API_URL" envDefault:"https://api.customeros.ai" validate:"required"`
	//}
	//Ai struct {
	//	ApiPath string `env:"AI_API_PATH,required" envDefault:"N/A"`
	//	ApiKey  string `env:"AI_API_KEY,required" envDefault:"N/A"`
	//}
	//OpenAi struct {
	//	Organization string `env:"OPENAI_ORGANIZATION,required" envDefault:""`
	//}
	//Anthropic struct {
	//	IndustryLookupPrompt1    string `env:"ANTHROPIC_INDUSTRY_LOOKUP_PROMPT,required" envDefault:"With next Global Industry Classification Standard (GICS) valid values: (Aerospace & Defense,Air Freight & Logistics,Automobile Components,Automobiles,Banks,Beverages,Biotechnology,Broadline Retail,Building Products,Capital Markets,Chemicals,Commercial Services & Supplies,Communications Equipment,Construction & Engineering,Construction Materials,Consumer Finance,Consumer Staples Distribution & Retail,Containers & Packaging,Diversified Consumer Services,Diversified REITs,Diversified Telecommunication Services,Distributors,Electric Utilities,Electrical Equipment,Electronic Equipment,Instruments & Components,Energy Equipment & Services,Entertainment,Financial Services,Food Products,Gas Utilities,Ground Transportation,Health Care Equipment & Supplies,Health Care Providers & Services,Health Care REITs,Health Care Technology,Hotel & Resort REITs,Hotels,Restaurants & Leisure,Household Durables,Household Products,Independent Power and Renewable Electricity Producers,Industrial Conglomerates,Industrial REITs,Insurance,Interactive Media & Services,Internet Software & Services,IT Services,Leisure Products,Life Sciences Tools & Services,Machinery,Marine Transportation,Media,Metals & Mining,Mortgage Real Estate Investment Trusts (REITs),Multi-Utilities,Office REITs,Oil,Gas & Consumable Fuels,Paper & Forest Products,Passenger Airlines,Personal Products,Pharmaceuticals,Professional Services,Real Estate Management & Development,Residential REITs,Retail REITs,Semiconductors & Semiconductor Equipment,Software,Specialized REITs,Specialty Retail,Technology Hardware,Storage & Peripherals,Textiles,Apparel & Luxury Goods,Tobacco,Trading Companies & Distributors,Transportation Infrastructure,Water Utilities,Wireless Telecommunication Services), provide appropriate industry mapping for (%s) and if do not see obvious mapping, provide appropriate GICS value from the input list based on other companies providing similar services. Finally if cannot map return just single word: Unknown"`
	//	IndustryLookupPrompt2    string `env:"ANTHROPIC_INDUSTRY_LOOKUP_PROMPT,required" envDefault:"What GICS value from following list (Aerospace & Defense,Air Freight & Logistics,Automobile Components,Automobiles,Banks,Beverages,Biotechnology,Broadline Retail,Building Products,Capital Markets,Chemicals,Commercial Services & Supplies,Communications Equipment,Construction & Engineering,Construction Materials,Consumer Finance,Consumer Staples Distribution & Retail,Containers & Packaging,Diversified Consumer Services,Diversified REITs,Diversified Telecommunication Services,Distributors,Electric Utilities,Electrical Equipment,Electronic Equipment,Instruments & Components,Energy Equipment & Services,Entertainment,Financial Services,Food Products,Gas Utilities,Ground Transportation,Health Care Equipment & Supplies,Health Care Providers & Services,Health Care REITs,Health Care Technology,Hotel & Resort REITs,Hotels,Restaurants & Leisure,Household Durables,Household Products,Independent Power and Renewable Electricity Producers,Industrial Conglomerates,Industrial REITs,Insurance,Interactive Media & Services,Internet Software & Services,IT Services,Leisure Products,Life Sciences Tools & Services,Machinery,Marine Transportation,Media,Metals & Mining,Mortgage Real Estate Investment Trusts (REITs),Multi-Utilities,Office REITs,Oil,Gas & Consumable Fuels,Paper & Forest Products,Passenger Airlines,Personal Products,Pharmaceuticals,Professional Services,Real Estate Management & Development,Residential REITs,Retail REITs,Semiconductors & Semiconductor Equipment,Software,Specialized REITs,Specialty Retail,Technology Hardware,Storage & Peripherals,Textiles,Apparel & Luxury Goods,Tobacco,Trading Companies & Distributors,Transportation Infrastructure,Water Utilities,Wireless Telecommunication Services) is chosen in next statement. Strictly provide the value only: %s"`
	//	EmailSummaryPrompt       string `env:"ANTHROPIC_EMAIL_SUMMARY_PROMPT,required" envDefault:"Make a 120 characters summary for this html email: %v"`
	//	EmailActionsItemsPrompt  string `env:"ANTHROPIC_EMAIL_ACTIONS_ITEMS_PROMPT,required" envDefault:"Give me the action points to be taken for the email. The criticality for the action points should be at least medium severity. return response in jSON format, key - \"items\", value - array of strings. The email is: %v"`
	//	LocationEnrichmentPrompt string `env:"ANTHROPIC_LOCATION_ENRICHMENT_PROMPT,required" envDefault:"Given the address '%s', please provide a JSON representation of the Location object with all available information. Use the following structure, filling in as many fields as possible based on the given address. If a field cannot be determined, omit it from the JSON output. Strictly return only the JSON.\n\n{\n    \"country\": \"string\",\n    \"countryCodeA2\": \"string\",\n    \"countryCodeA3\": \"string\",\n    \"region\": \"string\",\n    \"locality\": \"string\",\n    \"address\": \"string\",\n    \"address2\": \"string\",\n    \"zip\": \"string\",\n    \"addressType\": \"string\",\n    \"houseNumber\": \"string\",\n    \"postalCode\": \"string\",\n    \"plusFour\": \"string\",\n    \"commercial\": boolean,\n    \"predirection\": \"string\",\n    \"district\": \"string\",\n    \"street\": \"string\",\n    \"latitude\": number,\n    \"longitude\": number,\n    \"timeZone\": \"string\",\n    \"utcOffset\": number\n}"`
	//}
	//Novu struct {
	//	ApiKey string `env:"NOVU_API_KEY,required" envDefault:"N/A"`
	//}
	//FileStoreApiConfig fsc.FileStoreApiConfig
}
