package constants

import "time"

const (
	ServiceName                 = "CUSTOMER-OS-WEBHOOKS"
	AppSourceCustomerOsWebhooks = "customer-os-webhooks"

	ComponentService = "service"
	// Deprecated, all methods with this constants should be deprecated
	ComponentNeo4jRepository = "neo4jRepository"

	RequestMaxBodySizeCommon                      = 1 * 1024 * 1024   // 1 MB
	RequestMaxBodySizeMessages                    = 10 * 1024 * 1024  // 10 MB
	RequestMaxTimeout                             = 300 * time.Second // 5 minutes
	MaxRetryCheckDataInNeo4jAfterEventRequest int = 7
)
