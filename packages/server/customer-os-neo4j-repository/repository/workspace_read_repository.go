package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type WorkspaceReadRepository interface {
	Get(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetByName(ctx context.Context, tenant, name string) (*dbtype.Node, error)
}

type workspaceReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewWorkspaceReadRepository(driver *neo4j.DriverWithContext, database string) WorkspaceReadRepository {
	return &workspaceReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *workspaceReadRepository) Get(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkspaceReadRepository.Get")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})--(w:Workspace) return w`)
	params := map[string]any{
		"tenant": tenant,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	return result.([]*dbtype.Node), nil
}

func (r *workspaceReadRepository) GetByName(ctx context.Context, tenant, name string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkspaceReadRepository.GetByName")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	if result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `MATCH (t:Tenant {name:$tenant})--(w:Workspace{name:$workspaceName}) return w`,
			map[string]any{
				"tenant":        tenant,
				"workspaceName": name,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		if err != nil && err.Error() == "Result contains no more records" {
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		return result.(*dbtype.Node), nil
	}
}
