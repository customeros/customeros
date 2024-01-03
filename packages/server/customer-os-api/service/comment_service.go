package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type CommentService interface {
	GetCommentsForIssues(ctx context.Context, issueIds []string) (*entity.CommentEntities, error)
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

func (s *commentService) GetCommentsForIssues(ctx context.Context, issueIds []string) (*entity.CommentEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentService.GetCommentsForIssues")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("issueIds", issueIds))

	comments, err := s.repositories.CommentRepository.GetAllForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	commentEntities := entity.CommentEntities{}
	for _, v := range comments {
		commentEntity := s.mapDbNodeToCommentEntity(*v.Node)
		commentEntity.DataloaderKey = v.LinkedNodeId
		commentEntities = append(commentEntities, *commentEntity)
	}
	span.LogFields(log.Int("result count", len(commentEntities)))
	return &commentEntities, nil
}

func (s *commentService) mapDbNodeToCommentEntity(node dbtype.Node) *entity.CommentEntity {
	props := utils.GetPropsFromNode(node)
	comment := entity.CommentEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &comment
}
