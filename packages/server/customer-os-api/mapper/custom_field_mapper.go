package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapCustomFieldInputsToEntities(inputs []*model.CustomFieldInput) *entity.CustomFieldEntities {
	if inputs == nil {
		return nil
	}
	var result entity.CustomFieldEntities
	for _, singleInput := range inputs {
		result = append(result, *MapCustomFieldInputToEntity(singleInput))
	}
	return &result
}

func MapCustomFieldInputToEntity(input *model.CustomFieldInput) *entity.CustomFieldEntity {
	customFieldEntity := entity.CustomFieldEntity{
		Id:            input.ID,
		Name:          utils.IfNotNilString(input.Name),
		Value:         input.Value,
		DataType:      input.Datatype.String(),
		TemplateId:    input.TemplateID,
		Source:        neo4jentity.DataSourceOpenline,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
	}
	customFieldEntity.AdjustValueByDatatype()
	return &customFieldEntity
}

func MapCustomFieldUpdateInputToEntity(input *model.CustomFieldUpdateInput) *entity.CustomFieldEntity {
	customFieldEntity := entity.CustomFieldEntity{
		Id:            utils.StringPtr(input.ID),
		Name:          input.Name,
		DataType:      input.Datatype.String(),
		Value:         input.Value,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
	}
	customFieldEntity.AdjustValueByDatatype()
	return &customFieldEntity
}

func MapEntitiesToCustomFields(customFieldEntities *entity.CustomFieldEntities) []*model.CustomField {
	var customFields []*model.CustomField
	for _, customFieldEntity := range *customFieldEntities {
		customFields = append(customFields, MapEntityToCustomField(&customFieldEntity))
	}
	return customFields
}

func MapEntityToCustomField(entity *entity.CustomFieldEntity) *model.CustomField {
	var datatype = model.CustomFieldDataType(entity.DataType)
	if !datatype.IsValid() {
		datatype = model.CustomFieldDataTypeText
	}
	return &model.CustomField{
		ID:        *entity.Id,
		Name:      entity.Name,
		Datatype:  datatype,
		Value:     entity.Value,
		Source:    MapDataSourceToModel(entity.Source),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
