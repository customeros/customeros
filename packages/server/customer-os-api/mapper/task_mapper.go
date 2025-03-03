package mapper

import (
	enummapper "github.com/customeros/customeros/packages/server/customer-os-api/mapper/enum"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapEntityToTask(entity *neo4jentity.TaskEntity) *model.Task {
	if entity == nil {
		return nil
	}
	return &model.Task{
		ID:          entity.Id,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		DueAt:       entity.DueAt,
		Subject:     &entity.Subject,
		Description: &entity.Description,
		Status:      enummapper.MapTaskStatusToModel(entity.Status),
	}
}

func MapEntitiesToTasks(entities []*neo4jentity.TaskEntity) []*model.Task {
	var issues []*model.Task
	for _, issueEntity := range entities {
		issues = append(issues, MapEntityToTask(issueEntity))
	}
	return issues
}
