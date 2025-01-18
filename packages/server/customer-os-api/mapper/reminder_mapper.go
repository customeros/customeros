package mapper

import (
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapEntityToReminder(entity *neo4jentity.ReminderEntity) *model.Reminder {
	if entity == nil {
		return nil
	}

	metadata := &model.Metadata{
		ID:          entity.Id,
		Created:     entity.CreatedAt,
		LastUpdated: entity.UpdatedAt,
	}

	return &model.Reminder{
		Metadata:  metadata,
		Content:   &entity.Content,
		DueDate:   &entity.DueDate,
		Dismissed: &entity.Dismissed,
	}
}
