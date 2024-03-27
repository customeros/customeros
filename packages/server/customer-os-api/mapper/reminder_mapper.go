package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToReminder(entity *neo4jentity.ReminderEntity) *model.Reminder {
	if entity == nil {
		return nil
	}

	metadata := &model.Metadata{
		ID:            entity.Id,
		Created:       entity.CreatedAt,
		LastUpdated:   entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
	}

	return &model.Reminder{
		Metadata:  metadata,
		Content:   &entity.Content,
		DueDate:   &entity.DueDate,
		Dismissed: &entity.Dismissed,
	}
}
