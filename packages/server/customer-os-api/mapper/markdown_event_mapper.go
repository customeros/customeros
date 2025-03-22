package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	enummapper "github.com/customeros/customeros/packages/server/customer-os-api/mapper/enum"
)

func MapEntityToMarkdownEvent(entity *neo4jentity.MarkdownEventEntity) *model.MarkdownEvent {
	markdownEvent := model.MarkdownEvent{
		Metadata: &model.Metadata{
			ID:          entity.Id,
			Created:     entity.CreatedAt,
			LastUpdated: entity.UpdatedAt,
			Source:      enummapper.MapDataSourceToModel(entity.Source),
			AppSource:   entity.AppSource,
		},
		Content: utils.StringPtr(entity.Content),
	}
	return &markdownEvent
}

func MapEntitiesToMarkdownEvents(entities *neo4jentity.MarkdownEventEntities) []*model.MarkdownEvent {
	var markdownEvents []*model.MarkdownEvent
	for _, markdownEventEntity := range *entities {
		markdownEvents = append(markdownEvents, MapEntityToMarkdownEvent(&markdownEventEntity))
	}
	return markdownEvents
}
