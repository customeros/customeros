package config

type CloudflareConfig struct {
	Url    string `env:"CLOUDFLARE_URL" envDefault:"https://api.cloudflare.com/client/v4" validate:"required"`
	ApiKey string `env:"CLOUDFLARE_API_KEY" `
	Email  string `env:"CLOUDFLARE_API_EMAIL"`
}

type OpenSRSConfig struct {
	Url      string `env:"OPENSRS_URL" envDefault:"https://admin.a.hostedemail.com"`
	ApiKey   string `env:"OPENSRS_API_KEY"`
	Username string `env:"OPENSRS_API_USERNAME"`
}

type PostmarkConfig struct {
	Url                         string `env:"POSTMARK_URL" envDefault:"https://api.postmarkapp.com"`
	AccountApiKey               string `env:"POSTMARK_ACCOUNT_API_KEY"`
	DefaultInboundStreamWebhook string `env:"POSTMARK_DEFAULT_INBOUND_STREAM_WEBHOOK"`
}

type IntegrationAppConfig struct {
	WorkspaceKey                    string `env:"INTEGRATION_APP_WORKSPACE_KEY"`
	WorkspaceSecret                 string `env:"INTEGRATION_APP_WORKSPACE_SECRET"`
	ApiTriggerUrlCreatePaymentLinks string `env:"INTEGRATION_APP_API_TRIGGER_URL_CREATE_PAYMENT_LINKS"`
}

type SlackConfig struct {
	ClientID                        string `env:"SLACK_CLIENT_ID"`
	ClientSecret                    string `env:"SLACK_CLIENT_SECRET"`
	NotifyNewTenantRegisteredHook   string `env:"SLACK_NOTIFY_NEW_TENANT_REGISTERED_WEBHOOK"`
	InternalAlertsRegisteredWebhook string `env:"SLACK_INTERNAL_ALERTS_REGISTERED_WEBHOOK" envDefault:""`
}

type StripeConfig struct {
	ApiKey string `env:"STRIPE_API_KEY" envDefault:"N/A"`
}

type NovuCofig struct {
	ApiKey      string `env:"NOVU_API_KEY" envDefault:""`
	FronteraUrl string `env:"NOVU_FRONTERA_URL" envDefault:""`
}

type AnthropicConfig struct {
	ApiPath string `env:"ANTHROPIC_API_PATH" envDefault:"https://api.anthropic.com/v1/messages"`
	ApiKey  string `env:"ANTHROPIC_API_KEY" envDefault:""`
	AnthropicPrompts
}

type AnthropicPrompts struct {
	IndustryLookupPrompt1    string `env:"ANTHROPIC_INDUSTRY_LOOKUP_PROMPT,required" envDefault:"With next Global Industry Classification Standard (GICS) valid values: (Aerospace & Defense,Air Freight & Logistics,Automobile Components,Automobiles,Banks,Beverages,Biotechnology,Broadline Retail,Building Products,Capital Markets,Chemicals,Commercial Services & Supplies,Communications Equipment,Construction & Engineering,Construction Materials,Consumer Finance,Consumer Staples Distribution & Retail,Containers & Packaging,Diversified Consumer Services,Diversified REITs,Diversified Telecommunication Services,Distributors,Electric Utilities,Electrical Equipment,Electronic Equipment,Instruments & Components,Energy Equipment & Services,Entertainment,Financial Services,Food Products,Gas Utilities,Ground Transportation,Health Care Equipment & Supplies,Health Care Providers & Services,Health Care REITs,Health Care Technology,Hotel & Resort REITs,Hotels,Restaurants & Leisure,Household Durables,Household Products,Independent Power and Renewable Electricity Producers,Industrial Conglomerates,Industrial REITs,Insurance,Interactive Media & Services,Internet Software & Services,IT Services,Leisure Products,Life Sciences Tools & Services,Machinery,Marine Transportation,Media,Metals & Mining,Mortgage Real Estate Investment Trusts (REITs),Multi-Utilities,Office REITs,Oil,Gas & Consumable Fuels,Paper & Forest Products,Passenger Airlines,Personal Products,Pharmaceuticals,Professional Services,Real Estate Management & Development,Residential REITs,Retail REITs,Semiconductors & Semiconductor Equipment,Software,Specialized REITs,Specialty Retail,Technology Hardware,Storage & Peripherals,Textiles,Apparel & Luxury Goods,Tobacco,Trading Companies & Distributors,Transportation Infrastructure,Water Utilities,Wireless Telecommunication Services), provide appropriate industry mapping for (%s) and if do not see obvious mapping, provide appropriate GICS value from the input list based on other companies providing similar services. Finally if cannot map return just single word: Unknown"`
	IndustryLookupPrompt2    string `env:"ANTHROPIC_INDUSTRY_LOOKUP_PROMPT,required" envDefault:"What GICS value from following list (Aerospace & Defense,Air Freight & Logistics,Automobile Components,Automobiles,Banks,Beverages,Biotechnology,Broadline Retail,Building Products,Capital Markets,Chemicals,Commercial Services & Supplies,Communications Equipment,Construction & Engineering,Construction Materials,Consumer Finance,Consumer Staples Distribution & Retail,Containers & Packaging,Diversified Consumer Services,Diversified REITs,Diversified Telecommunication Services,Distributors,Electric Utilities,Electrical Equipment,Electronic Equipment,Instruments & Components,Energy Equipment & Services,Entertainment,Financial Services,Food Products,Gas Utilities,Ground Transportation,Health Care Equipment & Supplies,Health Care Providers & Services,Health Care REITs,Health Care Technology,Hotel & Resort REITs,Hotels,Restaurants & Leisure,Household Durables,Household Products,Independent Power and Renewable Electricity Producers,Industrial Conglomerates,Industrial REITs,Insurance,Interactive Media & Services,Internet Software & Services,IT Services,Leisure Products,Life Sciences Tools & Services,Machinery,Marine Transportation,Media,Metals & Mining,Mortgage Real Estate Investment Trusts (REITs),Multi-Utilities,Office REITs,Oil,Gas & Consumable Fuels,Paper & Forest Products,Passenger Airlines,Personal Products,Pharmaceuticals,Professional Services,Real Estate Management & Development,Residential REITs,Retail REITs,Semiconductors & Semiconductor Equipment,Software,Specialized REITs,Specialty Retail,Technology Hardware,Storage & Peripherals,Textiles,Apparel & Luxury Goods,Tobacco,Trading Companies & Distributors,Transportation Infrastructure,Water Utilities,Wireless Telecommunication Services) is chosen in next statement. Strictly provide the value only: %s"`
	EmailSummaryPrompt       string `env:"ANTHROPIC_EMAIL_SUMMARY_PROMPT,required" envDefault:"Make a 120 characters summary for this html email: %v"`
	EmailActionsItemsPrompt  string `env:"ANTHROPIC_EMAIL_ACTIONS_ITEMS_PROMPT,required" envDefault:"Give me the action points to be taken for the email. The criticality for the action points should be at least medium severity. return response in jSON format, key - \"items\", value - array of strings. The email is: %v"`
	LocationEnrichmentPrompt string `env:"ANTHROPIC_LOCATION_ENRICHMENT_PROMPT,required" envDefault:"Given the address '%s', please provide a JSON representation of the Location object with all available information. Use the following structure, filling in as many fields as possible based on the given address. If a field cannot be determined, omit it from the JSON output. Strictly return only the JSON.\n\n{\n    \"country\": \"string\",\n    \"countryCodeA2\": \"string\",\n    \"countryCodeA3\": \"string\",\n    \"region\": \"string\",\n    \"locality\": \"string\",\n    \"address\": \"string\",\n    \"address2\": \"string\",\n    \"zip\": \"string\",\n    \"addressType\": \"string\",\n    \"houseNumber\": \"string\",\n    \"postalCode\": \"string\",\n    \"plusFour\": \"string\",\n    \"commercial\": boolean,\n    \"predirection\": \"string\",\n    \"district\": \"string\",\n    \"street\": \"string\",\n    \"latitude\": number,\n    \"longitude\": number,\n    \"timeZone\": \"string\",\n    \"utcOffset\": number\n}"`
}
