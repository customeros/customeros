package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapContactGroupInputToEntity(input model.ContactGroupInput) *entity.ContactGroupEntity {
	return &entity.ContactGroupEntity{
		Name:          input.Name,
		Source:        entity.DataSourceOpenline,
		SourceOfTruth: entity.DataSourceOpenline,
	}
}

func MapContactGroupUpdateInputToEntity(input model.ContactGroupUpdateInput) *entity.ContactGroupEntity {
	return &entity.ContactGroupEntity{
		Id:            input.ID,
		Name:          input.Name,
		SourceOfTruth: entity.DataSourceOpenline,
	}
}

func MapEntityToContactGroup(entity *entity.ContactGroupEntity) *model.ContactGroup {
	return &model.ContactGroup{
		ID:        entity.Id,
		Name:      entity.Name,
		Source:    MapDataSourceToModel(entity.Source),
		CreatedAt: entity.CreatedAt,
	}
}

func MapEntitiesToContactGroups(entities *entity.ContactGroupEntities) []*model.ContactGroup {
	var contactGroups []*model.ContactGroup
	for _, contactGroupEntity := range *entities {
		contactGroups = append(contactGroups, MapEntityToContactGroup(&contactGroupEntity))
	}
	return contactGroups
}
