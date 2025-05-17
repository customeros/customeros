package api_comment

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"golang.org/x/net/context"

	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
)

type commentService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewCommentService(log logger.Logger, repositories *repository.Repositories) cosapi_interfaces.CommentService {
	return &commentService{
		log:          log,
		repositories: repositories,
	}
}

func (s *commentService) GetCommentsForIssues(ctx context.Context, issueIds []string) (*neo4jentity.CommentEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "CommentService.GetCommentsForIssues")
	defer spans.Finish()
	spans.LogObjectAsJson("issueIds", issueIds)

	comments, err := s.repositories.Neo4jRepositories.CommentReadRepository.GetAllForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	commentEntities := neo4jentity.CommentEntities{}
	for _, v := range comments {
		commentEntity := neo4jmapper.MapDbNodeToCommentEntity(v.Node)
		commentEntity.DataloaderKey = v.LinkedNodeId
		commentEntities = append(commentEntities, *commentEntity)
	}
	spans.LogKV("result count", len(commentEntities))
	return &commentEntities, nil
}
