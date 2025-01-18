package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	mapper "github.com/customeros/customeros/packages/server/customer-os-api/mapper/enum"
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
		WorkspaceLogo:        utils.StringPtrNillable(entity.WorkspaceLogo),
		WorkspaceName:        utils.StringPtrNillable(entity.WorkspaceName),
	}
}

func MapEntitiesToTenantSettingsOpportunityStages(entities []*postgresEntity.TenantSettingsOpportunityStage) []*model.TenantSettingsOpportunityStageConfiguration {
	list := make([]*model.TenantSettingsOpportunityStageConfiguration, 0)
	for _, entity := range entities {
		list = append(list, &model.TenantSettingsOpportunityStageConfiguration{
			ID:             entity.ID,
			Value:          entity.Value,
			Label:          entity.Label,
			Order:          entity.Order,
			Visible:        entity.Visible,
			LikelihoodRate: entity.LikelihoodRate,
		})
	}
	return list
}
