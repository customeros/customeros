package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LogEntryRepository interface {
	GetById(ctx context.Context, logEntryId string) (*dbtype.Node, error)
}

type logEntryRepository struct {
	driver *neo4j.DriverWithContext
}

func NewLogEntryRepository(driver *neo4j.DriverWithContext) LogEntryRepository {
	return &logEntryRepository{
		driver: driver,
	}
}

func (r *logEntryRepository) GetById(parentCtx context.Context, logEntryId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "LogEntryRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("logEntryId", logEntryId))

	query := fmt.Sprintf(`MATCH (l:LogEntry_%s {id:$id}) return l`, common.GetTenantFromContext(ctx))
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	if result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"id": logEntryId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}
