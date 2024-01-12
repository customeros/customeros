package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntityToRenewalRecord(renewalRecordEntity *entity.RenewalsRecordEntity) *model.RenewalRecord {
	if renewalRecordEntity == nil {
		return nil
	}
	organization := entity.OrganizationEntity{}
	orgModel := model.Organization{}
	contract := entity.ContractEntity{}
	contractModel := model.Contract{}
	opportunity := entity.OpportunityEntity{}
	opportunityModel := model.Opportunity{}

	if renewalRecordEntity.Organization != organization {
		orgModel = *MapEntityToOrganization(&renewalRecordEntity.Organization)
	}
	if renewalRecordEntity.Contract != contract {
		contractModel = *MapEntityToContract(&renewalRecordEntity.Contract)
	}
	if renewalRecordEntity.Opportunity != opportunity {
		opportunityModel = *MapEntityToOpportunity(&renewalRecordEntity.Opportunity)
	}
	return &model.RenewalRecord{
		Organization: &orgModel,
		Contract:     &contractModel,
		Opportunity:  &opportunityModel,
	}
}

func MapEntitiesToRenewalRecords(renewalRecordEntities *entity.RenewalsRecordEntities) []*model.RenewalRecord {
	var renewalRecords []*model.RenewalRecord
	for _, renewalRecordEntity := range *renewalRecordEntities {
		renewalRecords = append(renewalRecords, MapEntityToRenewalRecord(&renewalRecordEntity))
	}
	return renewalRecords
}
