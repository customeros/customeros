package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapContactTypeInputToEntity(input model.ContactTypeInput) *entity.ContactTypeEntity {
	contactTypeEntity := entity.ContactTypeEntity{
		Name: input.Name,
	}
	return &contactTypeEntity
}

func MapContactTypeUpdateInputToEntity(input model.ContactTypeUpdateInput) *entity.ContactTypeEntity {
	contactTypeEntity := entity.ContactTypeEntity{
		Id:   input.ID,
		Name: input.Name,
	}
	return &contactTypeEntity
}

func MapEntityToContactType(entity *entity.ContactTypeEntity) *model.ContactType {
	return &model.ContactType{
		ID:   entity.Id,
		Name: entity.Name,
	}
}

func MapEntitiesToContactTypes(entities *entity.ContactTypeEntities) []*model.ContactType {
	var contactTypes []*model.ContactType
	for _, contactTypeEntity := range *entities {
		contactTypes = append(contactTypes, MapEntityToContactType(&contactTypeEntity))
	}
	return contactTypes
}
