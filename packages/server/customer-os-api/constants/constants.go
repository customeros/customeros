package constants

const (
	ServiceName              = "CUSTOMER-OS-API"
	AppSourceCustomerOsApi   = "customer-os-api"
	ComponentResolver        = "resolver"
	ComponentService         = "service"
	ComponentNeo4jRepository = "neo4jRepository"

	UrlCustomerOsApi = "https://customeros.ai"
	UrlInvoices      = UrlCustomerOsApi + "/invoices"

	MaxRetriesCheckDataInNeo4jAfterEventRequest int = 5
)
