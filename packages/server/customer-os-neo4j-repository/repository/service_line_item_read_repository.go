package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type ServiceLineItemReadRepository interface {
	GetServiceLineItemById(ctx context.Context, tenant, serviceLineItemId string) (*dbtype.Node, error)
	GetAllForContract(ctx context.Context, tenant, contractId string) ([]*neo4j.Node, error)
	GetLatestServiceLineItemByParentId(ctx context.Context, tenant, serviceLineItemParentId string, beforeDate time.Time) (*dbtype.Node, error)
}

type serviceLineItemReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewServiceLineItemReadRepository(driver *neo4j.DriverWithContext, database string) ServiceLineItemReadRepository {
	return &serviceLineItemReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *serviceLineItemReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *serviceLineItemReadRepository) GetAllForContract(ctx context.Context, tenant, contractId string) ([]*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemReadRepository.GetAllForContract")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("contractId", contractId))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})-[:HAS_SERVICE]->(sli:ServiceLineItem)
							WHERE sli:ServiceLineItem_%s
							RETURN sli`, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*neo4j.Node))))
	return result.([]*neo4j.Node), nil
}

func (r *serviceLineItemReadRepository) GetServiceLineItemById(ctx context.Context, tenant, serviceLineItemId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemReadRepository.GetServiceLineItemById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, serviceLineItemId)

	cypher := fmt.Sprintf(`MATCH (sli:ServiceLineItem {id:$id}) WHERE sli:ServiceLineItem_%s RETURN sli`, tenant)
	params := map[string]any{
		"id": serviceLineItemId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

func (r *serviceLineItemReadRepository) GetLatestServiceLineItemByParentId(ctx context.Context, tenant, serviceLineItemParentId string, beforeDate time.Time) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemReadRepository.GetLatestServiceLineItemByParentId")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("serviceLineItemParentId", serviceLineItemParentId), log.Object("beforeDate", beforeDate))

	cypher := `MATCH (sli:ServiceLineItem {parentId:$parentId}) WHERE sli.startedAt < $before RETURN sli ORDER BY sli.startedAt DESC LIMIT 1`
	params := map[string]any{
		"tenant":   tenant,
		"parentId": serviceLineItemParentId,
		"before":   beforeDate.Add(time.Millisecond * 1),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}
