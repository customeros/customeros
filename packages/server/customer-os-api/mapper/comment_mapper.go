package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapEntityToComment(entity *neo4jentity.CommentEntity) *model.Comment {
	comment := model.Comment{
		ID:            entity.Id,
		Content:       utils.StringPtr(entity.Content),
		ContentType:   utils.StringPtr(entity.ContentType),
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
	}
	return &comment
}

func MapEntitiesToComments(entities *neo4jentity.CommentEntities) []*model.Comment {
	var comments []*model.Comment
	for _, commentEntity := range *entities {
		comments = append(comments, MapEntityToComment(&commentEntity))
	}
	return comments
}
