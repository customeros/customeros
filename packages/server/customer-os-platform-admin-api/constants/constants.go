package constants

import "time"

const (
	ServiceName                         = "CUSTOMER-OS-PLATFORM-ADMIN-API"
	AppSourceCustomerOsPlatformAdminApi = "customer-os-platform-admin-api"

	ComponentService         = "service"
	ComponentNeo4jRepository = "neo4jRepository"
	ComponentRepository      = "repository"

	RequestMaxBodySizeCommon                      = 1 * 1024 * 1024   // 1 MB
	RequestMaxBodySizeMessages                    = 10 * 1024 * 1024  // 10 MB
	RequestMaxTimeout                             = 300 * time.Second // 5 minutes
	MaxRetryCheckDataInNeo4jAfterEventRequest int = 10
	TimeoutIntervalMs                             = 100
)
