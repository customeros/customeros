package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

func MapEntityToTenantSettings(entity *neo4jentity.TenantSettingsEntity) *model.TenantSettings {
	if entity == nil {
		return nil
	}
	return &model.TenantSettings{
		BillingEnabled:       entity.InvoicingEnabled,
		LogoRepositoryFileID: utils.StringPtrNillable(entity.LogoRepositoryFileId),
		BaseCurrency:         utils.ToPtr(mapper.MapCurrencyToModel(entity.BaseCurrency)),
		LogoURL:              entity.LogoRepositoryFileId,
	}
}

func MapEntitiesToTenantSettingsOpportunityStages(entities []*postgresEntity.TenantSettingsOpportunityStage) []*model.TenantSettingsOpportunityStageConfiguration {
	list := make([]*model.TenantSettingsOpportunityStageConfiguration, 0)
	for _, entity := range entities {
		list = append(list, &model.TenantSettingsOpportunityStageConfiguration{
			ID:      entity.ID,
			Value:   entity.Value,
			Label:   entity.Label,
			Order:   entity.Order,
			Visible: entity.Visible,
		})
	}
	return list
}
