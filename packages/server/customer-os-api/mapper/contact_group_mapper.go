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

func MapEntityToContactGroup(contactGroup *entity.ContactGroupEntity) *model.ContactGroup {
	return &model.ContactGroup{
		ID:     contactGroup.Id,
		Name:   contactGroup.Name,
		Source: MapDataSourceToModel(contactGroup.Source),
	}
}

func MapEntitiesToContactGroups(contactGroupEntities *entity.ContactGroupEntities) []*model.ContactGroup {
	var contactGroups []*model.ContactGroup
	for _, contactGroupEntity := range *contactGroupEntities {
		contactGroups = append(contactGroups, MapEntityToContactGroup(&contactGroupEntity))
	}
	return contactGroups
}
