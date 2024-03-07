package entity

import neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

type RenewalsRecordEntity struct {
	Organization OrganizationEntity
	Contract     neo4jentity.ContractEntity
	Opportunity  OpportunityEntity
}

type RenewalsRecordEntities []RenewalsRecordEntity
