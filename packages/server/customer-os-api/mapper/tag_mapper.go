package mapper

import (
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapTagInputToEntity(input model.TagInput) *entity.TagEntity {
	tagEntity := entity.TagEntity{
		Name:          input.Name,
		Source:        neo4jentity.DataSourceOpenline,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &tagEntity
}

func MapTagUpdateInputToEntity(input model.TagUpdateInput) *entity.TagEntity {
	tagEntity := entity.TagEntity{
		Id:   input.ID,
		Name: input.Name,
	}
	return &tagEntity
}

func MapEntityToTag(entity entity.TagEntity) *model.Tag {
	return &model.Tag{
		ID:        entity.Id,
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Source:    MapDataSourceToModel(entity.Source),
		AppSource: entity.AppSource,
	}
}

func MapEntitiesToTags(entities *entity.TagEntities) []*model.Tag {
	var tags []*model.Tag
	for _, tagEntity := range *entities {
		tags = append(tags, MapEntityToTag(tagEntity))
	}
	return tags
}
