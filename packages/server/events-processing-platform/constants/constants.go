package constants

const (
	AppSourceEventProcessingPlatform = "event-processing-platform"
	AppSourceIntegrationApp          = "integration.app"
	AppSourceSyncCustomerOsData      = "sync-customer-os-data"

	ComponentNeo4jRepository = "neo4jRepository"
	ComponentService         = "service"

	SourceOpenline = "openline"

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

	WorkerID     = "workerID"
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

	EsInternalStreamPrefix = "$"
	EsAll                  = "$all"
	StreamTempPrefix       = "temp"

	StreamMetadataMaxCount              = 200
	StreamMetadataMaxAgeSeconds         = 7 * 24 * 60 * 60  // 7 days
	StreamMetadataMaxAgeSecondsExtended = 30 * 24 * 60 * 60 // 30 days
)
