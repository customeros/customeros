package constants

import "time"

const (
	ServiceName                 = "CUSTOMER-OS-WEBHOOKS"
	AppSourceCustomerOsWebhooks = "customer-os-webhooks"

	ComponentService         = "service"
	ComponentNeo4jRepository = "neo4jRepository"
	ComponentRepository      = "repository"

	RequestMaxBodySizeCommon                      = 1 * 1024 * 1024  // 1 MB
	RequestMaxTimeout                             = 30 * time.Second // 60 seconds
	MaxRetryCheckDataInNeo4jAfterEventRequest int = 10
	TimeoutIntervalMs                             = 100
)
