package cosapi_interfaces

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"

type CacheService interface {
	InitCache()
	GetStates() []*model.GCliItem
}
