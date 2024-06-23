package mapper

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntitiesToDomainNames(entities *neo4jentity.DomainEntities) []string {
	var domains []string
	for _, domainEntity := range *entities {
		domains = append(domains, domainEntity.Domain)
	}
	return domains
}
