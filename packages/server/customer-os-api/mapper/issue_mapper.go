package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapEntityToIssue(entity *neo4jentity.IssueEntity) *model.Issue {
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

func MapEntitiesToIssues(entities []*neo4jentity.IssueEntity) []*model.Issue {
	var issues []*model.Issue
	for _, issueEntity := range entities {
		issues = append(issues, MapEntityToIssue(issueEntity))
	}
	return issues
}
