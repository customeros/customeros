package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
)

func MapEntitiesToDomainNames(entities *entity.DomainEntities) []string {
	var domains []string
	for _, domainEntity := range *entities {
		domains = append(domains, domainEntity.Domain)
	}
	return domains
}
