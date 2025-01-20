package config

type AnthropicConfig struct {
	ApiPath string `env:"ANTHROPIC_API_PATH" envDefault:"https://api.anthropic.com/v1/messages"`
	ApiKey  string `env:"ANTHROPIC_API_KEY"`
	Prompts AnthropicPrompts
}

type AnthropicPrompts struct {
	EmailSummaryPrompt       string `env:"ANTHROPIC_EMAIL_SUMMARY_PROMPT,required" envDefault:"Make a 120 characters summary for this html email: %v"`
	EmailActionsItemsPrompt  string `env:"ANTHROPIC_EMAIL_ACTIONS_ITEMS_PROMPT,required" envDefault:"Give me the action points to be taken for the email. The criticality for the action points should be at least medium severity. return response in jSON format, key - \"items\", value - array of strings. The email is: %v"`
	LocationEnrichmentPrompt string `env:"ANTHROPIC_LOCATION_ENRICHMENT_PROMPT,required" envDefault:"Given the address ‘%s’, please provide a JSON representation of the Location object with all available information. Use the following structure, filling in as many fields as possible based on the given address. If a field cannot be determined, omit it from the JSON output. Strictly return only the JSON. { \"country\": \"string\", \"countryCodeA2\": \"string\", \"countryCodeA3\": \"string\", \"region\": \"string\", \"locality\": \"string\", \"address\": \"string\", \"address2\": \"string\", \"zip\": \"string\", \"addressType\": \"string\", \"houseNumber\": \"string\", \"postalCode\": \"string\", \"plusFour\": \"string\", \"commercial\": boolean, \"predirection\": \"string\", \"district\": \"string\", \"street\": \"string\", \"latitude\": number, \"longitude\": number, \"timeZone\": \"string\", \"utcOffset\": number }" `
}
