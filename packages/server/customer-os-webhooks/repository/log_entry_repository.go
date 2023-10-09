package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LogEntryRepository interface {
	GetById(ctx context.Context, tenant, logEntryId string) (*dbtype.Node, error)
	GetMatchedLogEntryId(ctx context.Context, tenant, externalSystem, externalId string) (string, error)
}

type logEntryRepository struct {
	driver *neo4j.DriverWithContext
}

func NewLogEntryRepository(driver *neo4j.DriverWithContext) LogEntryRepository {
	return &logEntryRepository{
		driver: driver,
	}
}

func (r *logEntryRepository) GetById(parentCtx context.Context, tenant, logEntryId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "LogEntryRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("logEntryId", logEntryId))

	query := fmt.Sprintf(`MATCH (log:LogEntry_%s {id:$logEntryId}) RETURN log`, tenant)
	params := map[string]any{
		"logEntryId": logEntryId,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}

func (r *logEntryRepository) GetMatchedLogEntryId(ctx context.Context, tenant, externalSystem, externalId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryRepository.GetMatchedLogEntryId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", externalSystem), log.String("externalId", externalId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (e)<-[:IS_LINKED_WITH {externalId:$logEntryExternalId}]-(l:LogEntry)
				WITH l WHERE l is not null
				return l.id order by l.createdAt limit 1`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":             tenant,
				"externalSystem":     externalSystem,
				"logEntryExternalId": externalId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	noteIDs := dbRecords.([]*db.Record)
	if len(noteIDs) > 0 {
		return noteIDs[0].Values[0].(string), nil
	}
	return "", nil
}
