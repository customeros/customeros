package neo4j_repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TaskReadRepository interface {
	GetById(ctx context.Context, tenant, taskId string) (*dbtype.Node, error)
}

type taskReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func (r *taskReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func NewTaskReadRepository(driver *neo4j.DriverWithContext, database string) TaskReadRepository {
	return &taskReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *taskReadRepository) GetById(ctx context.Context, tenant, taskId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TaskReadRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	tracing.TagEntity(span, taskId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:TASK_BELONGS_TO_TENANT]-(tsk:Task {id:$taskId})
				RETURN tsk`
	params := map[string]any{
		"tenant": tenant,
		"taskId": taskId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}
