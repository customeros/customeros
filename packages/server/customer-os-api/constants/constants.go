package constants

const (
	ServiceName                = "CUSTOMER-OS-API"
	AppSourceCustomerOsApi     = "customer-os-api"
	AppSourceCustomerOsApiRest = "customer-os-api/rest"

	UrlCustomerOsApi = "https://customeros.ai"
	UrlInvoices      = UrlCustomerOsApi + "/invoices"

	MaxRetriesCheckDataInNeo4jAfterEventRequest int = 5
	FileStoreFileDownloadUrlTemplate                = UrlCustomerOsApi + "/files/v1/files/%s/download"
)
