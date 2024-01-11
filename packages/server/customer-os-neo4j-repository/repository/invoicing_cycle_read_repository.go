package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type InvoicingCycleReadRepository interface {
	GetInvoicingCycle(ctx context.Context, tenant string) (*dbtype.Node, error)
}

type invocingCycleReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInvoicingCycleReadRepository(driver *neo4j.DriverWithContext, database string) InvoicingCycleReadRepository {
	return &invocingCycleReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *invocingCycleReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *invocingCycleReadRepository) GetInvoicingCycle(ctx context.Context, tenant string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanReadRepository.GetMasterPlanById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICING_CYCLE_BELONGS_TO_TENANT]-(ic:InvoicingCycle) RETURN ic`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)

	})
	if err != nil {
		span.LogFields(log.Bool("result.found", false))
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}
