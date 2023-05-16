package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapPersonInputToEntity(input *model.PersonInput) *entity.PersonEntity {
	if input == nil {
		return nil
	}
	personEntity := entity.PersonEntity{
		IdentityId:    input.IdentityID,
		Email:         input.Email,
		Provider:      input.Provider,
		CreatedAt:     utils.Now(),
		UpdatedAt:     utils.Now(),
		Source:        entity.DataSourceOpenline,
		SourceOfTruth: entity.DataSourceOpenline,
		AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &personEntity
}

func MapPersonUpdateToEntity(id string, input *model.PersonUpdate) *entity.PersonEntity {
	if input == nil {
		return nil
	}
	personEntity := entity.PersonEntity{
		Id:            id,
		IdentityId:    input.IdentityID,
		UpdatedAt:     utils.Now(),
		SourceOfTruth: entity.DataSourceOpenline,
		AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &personEntity
}

func MapEntityToPerson(entity *entity.PersonEntity) *model.Person {
	return &model.Person{
		ID:            entity.Id,
		IdentityID:    entity.IdentityId,
		Email:         entity.Email,
		Provider:      entity.Provider,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
	}
}
