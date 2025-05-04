package neo4j_repository

import (
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	"golang.org/x/net/context"
)

type LogEntryReadRepository interface {
	GetById(ctx context.Context, tenant, id string) (*dbtype.Node, error)
}

type logEntryReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewLogEntryReadRepository(driver *neo4j.DriverWithContext, database string) LogEntryReadRepository {
	return &logEntryReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *logEntryReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *logEntryReadRepository) GetById(ctx context.Context, tenant, id string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "LogEntryReadRepository.GetById")
	defer spans.Finish()

	spans.TagEntity(id)

	cypher := fmt.Sprintf(`MATCH (l:LogEntry {id:$id}) WHERE l:LogEntry_%s return l`, tenant)
	params := map[string]any{
		"id": id,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		spans.LogKV("result.found", false)
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}
