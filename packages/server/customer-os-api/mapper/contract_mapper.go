package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToContract(entity *entity.ContractEntity) *model.Contract {
	if entity == nil {
		return nil
	}
	return &model.Contract{
		ID:               entity.ID,
		Name:             entity.Name,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
		Source:           MapDataSourceToModel(entity.Source),
		SourceOfTruth:    MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:        entity.AppSource,
		Status:           MapContractStatusToModel(entity.ContractStatus),
		RenewalCycle:     MapContractRenewalCycleToModel(entity.ContractRenewalCycle),
		ServiceStartedAt: entity.ServiceStartedAt,
		SignedAt:         entity.SignedAt,
		EndedAt:          entity.EndedAt,
		ContractURL:      utils.StringPtrNillable(entity.ContractUrl),
	}
}
