package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapEntityDefinitionInputToEntity(input model.EntityDefinitionInput) *entity.EntityDefinitionEntity {
	definitionEntity := entity.EntityDefinitionEntity{
		Name: input.Name,
	}
	if input.Extends != nil {
		extends := input.Extends.String()
		definitionEntity.Extends = &extends
	}
	return &definitionEntity
}

func MapEntityToEntityDefinition(entity *entity.EntityDefinitionEntity) *model.EntityDefinition {
	output := model.EntityDefinition{
		ID:      entity.Id,
		Name:    entity.Name,
		Version: int(entity.Version),
	}
	if entity.Extends != nil {
		extends := model.EntityDefinitionExtension(*entity.Extends)
		if extends.IsValid() {
			output.Extends = &extends
		}
	}
	return &output
}
