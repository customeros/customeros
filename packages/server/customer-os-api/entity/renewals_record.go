package entity

type RenewalsRecordEntity struct {
	Organization OrganizationEntity
	Contract     ContractEntity
	Opportunity  OpportunityEntity
}

type RenewalsRecordEntities []RenewalsRecordEntity
