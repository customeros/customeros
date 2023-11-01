package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToComment(entity *entity.CommentEntity) *model.Comment {
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

func MapEntitiesToComments(entities *entity.CommentEntities) []*model.Comment {
	var comments []*model.Comment
	for _, commentEntity := range *entities {
		comments = append(comments, MapEntityToComment(&commentEntity))
	}
	return comments
}
