package cosapi_interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type CommentService interface {
	GetCommentsForIssues(ctx context.Context, issueIds []string) (*neo4j_entity.CommentEntities, error)
}
