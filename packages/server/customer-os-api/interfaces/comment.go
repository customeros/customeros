package cosapi_interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type CommentService interface {
	GetCommentsForIssues(ctx context.Context, issueIds []string) (*entity.CommentEntities, error)
}
