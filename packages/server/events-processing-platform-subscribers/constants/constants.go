package constants

const (
	AppSourceEventProcessingPlatformSubscribers = "event-processing-platform-subscribers"
	AppSourceCustomerOsApi                      = "customer-os-api"

	ComponentSubscriptionGraph   = "subscriptionGraph"
	ComponentSubscriptionInvoice = "subscriptionInvoice"

	AggregateTypeOpportunity = "opportunity"

	SourceOpenline = "openline"

	AppBrandfetch = "brandfetch"
	AppScrapin    = "scrapin"
	AppEnrichment = "enrichment-api"

	PromptType_EmailSummary         = "EmailSummary"
	PromptType_EmailActionItems     = "EmailActionItems"
	PromptType_MapIndustry          = "MapIndustryToList"
	PromptType_ExtractIndustryValue = "ExtractIndustryValueFromAiResponse"
	PromptTypeExtractLocationValue  = "ExtractLocationValue"

	Anthropic         = "anthropic"
	OpenAI            = "openai"
	AnthropicApiModel = "claude-3-5-sonnet-20240620"
	//AnthropicApiModel = "claude-3-haiku-20240307"

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

	MaxRetriesCheckDataInNeo4j = 8

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

	WorkerID     = "workerID"
	GroupName    = "GroupName"
	StreamID     = "StreamID"
	EventID      = "EventID"
	EventType    = "EventType"
	EventNumber  = "EventNumber"
	CreatedDate  = "CreatedDate"
	UserMetadata = "UserMetadata"

	EsInternalStreamPrefix = "$"

	UrlFileStoreApi                     = "https://fs.customeros.ai/"
	UrlFileStoreFileDownloadUrlTemplate = UrlFileStoreApi + "file/%s/download"
)
