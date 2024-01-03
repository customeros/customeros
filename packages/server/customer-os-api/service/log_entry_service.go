package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LogEntryService interface {
	GetById(ctx context.Context, logEntryId string) (*entity.LogEntryEntity, error)
	mapDbNodeToLogEntryEntity(node *dbtype.Node) *entity.LogEntryEntity
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

func (s *logEntryService) GetById(ctx context.Context, logEntryId string) (*entity.LogEntryEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("logEntryId", logEntryId))

	logEntryDbNode, err := s.repositories.LogEntryRepository.GetById(ctx, logEntryId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return s.mapDbNodeToLogEntryEntity(logEntryDbNode), nil
}

func (s *logEntryService) mapDbNodeToLogEntryEntity(node *dbtype.Node) *entity.LogEntryEntity {
	props := utils.GetPropsFromNode(*node)
	logEntry := entity.LogEntryEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		StartedAt:     utils.GetTimePropOrEpochStart(props, "startedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &logEntry
}
