package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapTagInputToEntity(input model.TagInput) *neo4jentity.TagEntity {
	tagEntity := neo4jentity.TagEntity{
		Name:          input.Name,
		Source:        neo4jentity.DataSourceOpenline,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &tagEntity
}

func MapEntityToTag(entity *neo4jentity.TagEntity) *model.Tag {
	if entity == nil {
		return nil
	}
	return &model.Tag{
		Metadata: &model.Metadata{
			ID:            entity.Id,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        MapDataSourceToModel(entity.Source),
			SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
			AppSource:     entity.AppSource,
		},
		ID:        entity.Id,
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Source:    MapDataSourceToModel(entity.Source),
		AppSource: entity.AppSource,
	}
}

func MapEntitiesToTags(entities *neo4jentity.TagEntities) []*model.Tag {
	var tags []*model.Tag
	for _, tagEntity := range *entities {
		tags = append(tags, MapEntityToTag(&tagEntity))
	}
	return tags
}
