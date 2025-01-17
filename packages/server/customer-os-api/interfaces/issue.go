package cosapi_interfaces

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type IssueService interface {
	GetIssueSummaryByStatusForOrganization(ctx context.Context, organizationId string) (map[string]int64, error)
	GetById(ctx context.Context, issueId string) (*entity.IssueEntity, error)
	GetIssuesForInteractionEvents(ctx context.Context, ids []string) (*entity.IssueEntities, error)
	GetSubmitterParticipantsForIssues(ctx context.Context, ids []string) (*neo4jentity.IssueParticipants, error)
	GetReporterParticipantsForIssues(ctx context.Context, ids []string) (*neo4jentity.IssueParticipants, error)
	GetAssigneeParticipantsForIssues(ctx context.Context, ids []string) (*neo4jentity.IssueParticipants, error)
	GetFollowerParticipantsForIssues(ctx context.Context, ids []string) (*neo4jentity.IssueParticipants, error)
}
