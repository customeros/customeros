package constants

const (
	AppSourceEventProcessingPlatform = "event-processing-platform"
	AppSourceIntegrationApp          = "integration.app"
	AppSourceSyncCustomerOsData      = "sync-customer-os-data"

	ComponentNeo4jRepository   = "neo4jRepository"
	ComponentService           = "service"
	ComponentCommandHandler    = "commandHandler"
	ComponentSubscriptionGraph = "subscriptionGraph"

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

	RenewalLikelihood_Order_High   = 40
	RenewalLikelihood_Order_Medium = 30
	RenewalLikelihood_Order_Low    = 20
	RenewalLikelihood_Order_Zero   = 10

	OnboardingStatus_Order_NotStarted = 10
	OnboardingStatus_Order_Stuck      = 20
	OnboardingStatus_Order_Late       = 30
	OnboardingStatus_Order_OnTrack    = 40
	OnboardingStatus_Order_Done       = 50
	OnboardingStatus_Order_Successful = 60

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
