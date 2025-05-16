package repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type LogEntryRepository interface {
	// Deprecated
	GetById(ctx context.Context, tenant, logEntryId string) (*dbtype.Node, error)
	// Deprecated
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
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "LogEntryRepository.GetById")
	defer spans.Finish()
	spans.LogKV("logEntryId", logEntryId)

	query := fmt.Sprintf(`MATCH (log:LogEntry_%s {id:$logEntryId}) RETURN log`, tenant)
	params := map[string]any{
		"logEntryId": logEntryId,
	}
	spans.LogKV("query", query)
	spans.LogObjectAsJson("params", params)

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "LogEntryRepository.GetMatchedLogEntryId")
	defer spans.Finish()
	spans.LogKV("externalSystem", externalSystem)
	spans.LogKV("externalId", externalId)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (e)<-[:IS_LINKED_WITH {externalId:$logEntryExternalId}]-(l:LogEntry)
				WITH l WHERE l is not null
				return l.id order by l.createdAt limit 1`
	spans.LogKV("query", query)

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
		spans.TraceError(err)
		return "", err
	}
	noteIDs := dbRecords.([]*db.Record)
	if len(noteIDs) > 0 {
		return noteIDs[0].Values[0].(string), nil
	}
	return "", nil
}
