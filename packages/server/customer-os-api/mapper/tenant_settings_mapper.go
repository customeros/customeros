package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToTenantSettings(entity *neo4jentity.TenantSettingsEntity) *model.TenantSettings {
	if entity == nil {
		return nil
	}
	return &model.TenantSettings{
		InvoicingEnabled: entity.InvoicingEnabled,
		LogoURL:          entity.LogoUrl,
		DefaultCurrency:  utils.ToPtr(MapCurrencyToModel(entity.DefaultCurrency)),
	}
}
