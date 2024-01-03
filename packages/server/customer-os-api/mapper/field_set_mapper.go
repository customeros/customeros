package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapFieldSetInputsToEntities(inputs []*model.FieldSetInput) *entity.FieldSetEntities {
	if inputs == nil {
		return nil
	}
	var result entity.FieldSetEntities
	for _, singleInput := range inputs {
		result = append(result, *MapFieldSetInputToEntity(singleInput))
	}
	return &result
}

func MapFieldSetInputToEntity(input *model.FieldSetInput) *entity.FieldSetEntity {
	fieldSetEntity := entity.FieldSetEntity{
		Id:            input.ID,
		Name:          input.Name,
		TemplateId:    input.TemplateID,
		Source:        neo4jentity.DataSourceOpenline,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		CustomFields:  MapCustomFieldInputsToEntities(input.CustomFields),
	}
	return &fieldSetEntity
}

func MapFieldSetUpdateInputToEntity(input *model.FieldSetUpdateInput) *entity.FieldSetEntity {
	fieldSetEntity := entity.FieldSetEntity{
		Id:            utils.StringPtr(input.ID),
		Name:          input.Name,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
	}
	return &fieldSetEntity
}

func MapEntitiesToFieldSets(fieldSetEntities *entity.FieldSetEntities) []*model.FieldSet {
	var fieldSet []*model.FieldSet
	for _, fieldSetEntity := range *fieldSetEntities {
		fieldSet = append(fieldSet, MapEntityToFieldSet(&fieldSetEntity))
	}
	return fieldSet
}

func MapEntityToFieldSet(fieldSetEntity *entity.FieldSetEntity) *model.FieldSet {
	return &model.FieldSet{
		ID:        *fieldSetEntity.Id,
		Name:      fieldSetEntity.Name,
		CreatedAt: fieldSetEntity.CreatedAt,
		UpdatedAt: fieldSetEntity.UpdatedAt,
		Source:    MapDataSourceToModel(fieldSetEntity.Source),
	}
}
