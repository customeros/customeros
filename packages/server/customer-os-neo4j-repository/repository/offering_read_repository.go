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

type OfferingReadRepository interface {
	GetOfferings(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetOfferingById(ctx context.Context, tenant, id string) (*dbtype.Node, error)
}

type offeringReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOfferingReadRepository(driver *neo4j.DriverWithContext, database string) OfferingReadRepository {
	return &offeringReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *offeringReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *offeringReadRepository) GetOfferings(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OfferingReadRepository.GetOfferings")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:OFFERING_BELONGS_TO_TENANT]-(of:Offering)
			RETURN of ORDER BY of.createdAt ASC`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})

	if err != nil {
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	if result == nil {
		return nil, nil
	}
	return result.([]*dbtype.Node), nil
}

func (r *offeringReadRepository) GetOfferingById(ctx context.Context, tenant, id string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OfferingReadRepository.GetOfferingById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:OFFERING_BELONGS_TO_TENANT]-(of:Offering {id:$id})
			RETURN of`
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
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
		tracing.TraceErr(span, err)
		span.LogFields(log.Bool("result.found", false))
		return nil, err
	}

	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}
