package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToIssue(entity *entity.IssueEntity) *model.Issue {
	if entity == nil {
		return nil
	}
	return &model.Issue{
		ID:            entity.Id,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Subject:       utils.StringPtr(entity.Subject),
		Status:        entity.Status,
		IssueStatus:   entity.Status,
		Priority:      utils.StringPtr(entity.Priority),
		Description:   utils.StringPtr(entity.Description),
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
	}
}

func MapEntitiesToIssues(entities []*entity.IssueEntity) []*model.Issue {
	var issues []*model.Issue
	for _, issueEntity := range entities {
		issues = append(issues, MapEntityToIssue(issueEntity))
	}
	return issues
}
