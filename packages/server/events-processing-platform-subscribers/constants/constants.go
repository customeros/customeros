package constants

const (
	AppSourceEventProcessingPlatformSubscribers = "event-processing-platform-subscribers"

	ComponentSubscriptionGraph   = "subscriptionGraph"
	ComponentSubscriptionInvoice = "subscriptionInvoice"

	SourceOpenline = "openline"

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

	UrlCustomerOsApi                 = "https://customeros.ai"
	FileStoreFileDownloadUrlTemplate = UrlCustomerOsApi + "/files/v1/files/%s/download"
)
