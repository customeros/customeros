package constants

const (
	AppSourceEventProcessingPlatform = "event-processing-platform"

	ComponentNeo4jRepository = "neo4jRepository"
	ComponentService         = "service"
	ComponentCommandHandler  = "commandHandler"

	SourceOpenline  = "openline"
	SourceWebscrape = "webscrape"

	PromptType_EmailSummary                = "EmailSummary"
	PromptType_EmailActionItems            = "EmailActionItems"
	PromptType_MapIndustry                 = "MapIndustryToList"
	PromptType_ExtractIndustryValue        = "ExtractIndustryValueFromAiResponse"
	PromptType_WebscrapeCompanyPrompt      = "CompanyAnalysisFromWebsite"
	PromptType_WebscrapeExtractCompanyData = "ExtractCompanyDataFromAnalysis"

	Anthropic = "anthropic"
	OpenAI    = "openai"

	NodeLabel_Organization     = "Organization"
	NodeLabel_InteractionEvent = "InteractionEvent"
	NodeLabel_User             = "User"
	NodeLabel_Contact          = "Contact"
	NodeLabel_LogEntry         = "LogEntry"

	TenantKeyHeader = "X-OPENLINE-TENANT-KEY"
	ApiKeyHeader    = "X-Openline-API-KEY"

	Tcp = "tcp"

	GRPC     = "GRPC"
	SIZE     = "SIZE"
	URI      = "URI"
	STATUS   = "STATUS"
	HTTP     = "HTTP"
	ERROR    = "ERROR"
	METHOD   = "METHOD"
	METADATA = "METADATA"
	REQUEST  = "REQUEST"
	REPLY    = "REPLY"
	TIME     = "TIME"

	Topic        = "topic"
	Partition    = "partition"
	Message      = "message"
	WorkerID     = "workerID"
	Offset       = "offset"
	Time         = "time"
	GroupName    = "GroupName"
	StreamID     = "StreamID"
	EventID      = "EventID"
	EventType    = "EventType"
	EventNumber  = "EventNumber"
	CreatedDate  = "CreatedDate"
	UserMetadata = "UserMetadata"

	Validate        = "validate"
	FieldValidation = "field validation"
	RequiredHeaders = "required header"
	Base64          = "base64"
	Unmarshal       = "unmarshal"
	Uuid            = "uuid"
	Cookie          = "cookie"
	Token           = "token"
	Bcrypt          = "bcrypt"
	Redis           = "redis"

	EsAll = "$all"
)
