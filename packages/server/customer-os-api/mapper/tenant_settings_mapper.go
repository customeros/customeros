package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToTenantSettings(entity *neo4jentity.TenantSettingsEntity) *model.TenantSettings {
	if entity == nil {
		return nil
	}
	return &model.TenantSettings{
		BillingEnabled:       entity.InvoicingEnabled,
		LogoRepositoryFileID: utils.StringPtrNillable(entity.LogoRepositoryFileId),
		BaseCurrency:         utils.ToPtr(mapper.MapCurrencyToModel(entity.BaseCurrency)),
		OpportunityStages:    entity.OpportunityStages,
		LogoURL:              entity.LogoRepositoryFileId,
	}
}
