package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LogEntryService interface {
	GetById(ctx context.Context, logEntryId string) (*neo4jentity.LogEntryEntity, error)
}

type logEntryService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewLogEntryService(log logger.Logger, repositories *repository.Repositories) LogEntryService {
	return &logEntryService{
		log:          log,
		repositories: repositories,
	}
}

func (s *logEntryService) GetById(ctx context.Context, logEntryId string) (*neo4jentity.LogEntryEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("logEntryId", logEntryId))

	logEntryDbNode, err := s.repositories.Neo4jRepositories.LogEntryReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), logEntryId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return neo4jmapper.MapDbNodeToLogEntryEntity(logEntryDbNode), nil
}
