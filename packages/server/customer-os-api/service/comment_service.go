package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type CommentService interface {
	GetCommentsForIssues(ctx context.Context, issueIds []string) (*neo4jentity.CommentEntities, error)
}

type commentService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewCommentService(log logger.Logger, repositories *repository.Repositories) CommentService {
	return &commentService{
		log:          log,
		repositories: repositories,
	}
}

func (s *commentService) GetCommentsForIssues(ctx context.Context, issueIds []string) (*neo4jentity.CommentEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentService.GetCommentsForIssues")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("issueIds", issueIds))

	comments, err := s.repositories.Neo4jRepositories.CommentReadRepository.GetAllForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	commentEntities := neo4jentity.CommentEntities{}
	for _, v := range comments {
		commentEntity := neo4jmapper.MapDbNodeToCommentEntity(v.Node)
		commentEntity.DataloaderKey = v.LinkedNodeId
		commentEntities = append(commentEntities, *commentEntity)
	}
	span.LogFields(log.Int("result count", len(commentEntities)))
	return &commentEntities, nil
}
