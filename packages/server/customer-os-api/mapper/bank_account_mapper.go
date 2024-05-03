package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToBankAccount(entity *neo4jentity.BankAccountEntity) *model.BankAccount {
	if entity == nil {
		return nil
	}
	return &model.BankAccount{
		Metadata: &model.Metadata{
			ID:            entity.Id,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        MapDataSourceToModel(entity.Source),
			SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
			AppSource:     entity.AppSource,
		},
		BankName:            utils.StringPtr(entity.BankName),
		BankTransferEnabled: entity.BankTransferEnabled,
		AllowInternational:  entity.AllowInternational,
		Currency:            utils.ToPtr(mapper.MapCurrencyToModel(entity.Currency)),
		Iban:                utils.StringPtr(entity.Iban),
		Bic:                 utils.StringPtr(entity.Bic),
		SortCode:            utils.StringPtr(entity.SortCode),
		AccountNumber:       utils.StringPtr(entity.AccountNumber),
		RoutingNumber:       utils.StringPtr(entity.RoutingNumber),
		OtherDetails:        utils.StringPtr(entity.OtherDetails),
	}
}

func MapEntitiesToBankAccounts(entities *neo4jentity.BankAccountEntities) []*model.BankAccount {
	var models []*model.BankAccount
	for _, entity := range *entities {
		models = append(models, MapEntityToBankAccount(&entity))
	}
	return models
}
