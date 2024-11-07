package constants

const (
	AppSourceCustomerOsApi = "customer-os-api"
	AppSourceUserAdminApi  = "user-admin-api"
	AppSourceFileStoreApi  = "file-store-api"

	ComponentResolver        = "resolver"
	ComponentRest            = "rest"
	ComponentService         = "service"
	ComponentListener        = "listener"
	ComponentNeo4jRepository = "neo4jRepository"
	// Deprecated: Use tracing package instead
	ComponentPostgresRepository = "postgresRepository"

	PromptType_EmailSummary         = "EmailSummary"
	PromptType_EmailActionItems     = "EmailActionItems"
	PromptType_MapIndustry          = "MapIndustryToList"
	PromptType_ExtractIndustryValue = "ExtractIndustryValueFromAiResponse"
	PromptTypeExtractLocationValue  = "ExtractLocationValue"

	Anthropic         = "anthropic"
	OpenAI            = "openai"
	AnthropicApiModel = "claude-3-5-sonnet-20240620"
	//AnthropicApiModel = "claude-3-haiku-20240307"
)
