package entity

import neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

type RenewalsRecordEntity struct {
	Organization neo4jentity.OrganizationEntity
	Contract     neo4jentity.ContractEntity
	Opportunity  neo4jentity.OpportunityEntity
}

type RenewalsRecordEntities []RenewalsRecordEntity
