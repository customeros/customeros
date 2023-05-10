package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapDomainInputToEntity(input model.DomainInput) entity.DomainEntity {
	domainEntity := entity.DomainEntity{
		Domain:    input.Domain,
		Source:    entity.DataSourceOpenline,
		AppSource: utils.IfNotNilString(input.AppSource),
	}
	if len(domainEntity.AppSource) == 0 {
		domainEntity.AppSource = constants.AppSourceCustomerOsApi
	}
	return domainEntity
}

func MapEntitiesToDomainNames(entities *entity.DomainEntities) []string {
	var domains []string
	for _, domainEntity := range *entities {
		domains = append(domains, domainEntity.Domain)
	}
	return domains
}
