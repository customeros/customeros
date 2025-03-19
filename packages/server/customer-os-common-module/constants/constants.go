package constants

const (
	AppSourceCustomerOsApi     = "customer-os-api"
	AppSourceSyncEmail         = "sync-email"
	AppSourceUpkeeper          = "customer-os-data-upkeeper"
	AppSourceEventsSubscribers = "events-subscribers"

	ComponentResolver = "resolver"
	ComponentListener = "listener"

	UrlCustomerOsApi                 = "https://customeros.ai"
	FileStoreFileDownloadUrlTemplate = UrlCustomerOsApi + "/files/v1/files/%s/download"

	S3ImagesCDN = "https://customer-os-images.b-cdn.net/"
)
