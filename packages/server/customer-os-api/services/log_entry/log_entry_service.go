package api_log_entry

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
)

type logEntryService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewLogEntryService(log logger.Logger, repositories *repository.Repositories) cosapi_interfaces.LogEntryService {
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
