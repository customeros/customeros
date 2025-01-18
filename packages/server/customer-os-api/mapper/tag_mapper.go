package mapper

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
)

func MapTagInputToEntity(input model.TagInput) *neo4jentity.TagEntity {
	tagEntity := neo4jentity.TagEntity{
		Name:      input.Name,
		Source:    neo4jentity.DataSourceOpenline,
		ColorCode: utils.IfNotNilString(input.ColorCode),
	}
	if input.EntityType != nil {
		tagEntity.EntityType = commonmodel.DecodeEntityType(input.EntityType.String())
	}
	return &tagEntity
}

func MapEntityToTag(entity *neo4jentity.TagEntity) *model.Tag {
	if entity == nil {
		return nil
	}
	return &model.Tag{
		Metadata: &model.Metadata{
			ID:          entity.Id,
			Created:     entity.CreatedAt,
			LastUpdated: entity.UpdatedAt,
			AppSource:   entity.AppSource,
			Source:      MapDataSourceToModel(entity.Source),
		},
		Name:       entity.Name,
		ColorCode:  entity.ColorCode,
		EntityType: model.EntityType(entity.EntityType.String()),
		ID:         utils.StringPtr(entity.Id),
	}
}

func MapEntitiesToTags(entities *neo4jentity.TagEntities) []*model.Tag {
	var tags []*model.Tag
	for _, tagEntity := range *entities {
		tags = append(tags, MapEntityToTag(&tagEntity))
	}
	return tags
}
