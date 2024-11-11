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
	"time"
)

type ServiceLineItemReadRepository interface {
	GetServiceLineItemById(ctx context.Context, tenant, serviceLineItemId string) (*dbtype.Node, error)
	GetServiceLineItemsByParentId(ctx context.Context, tenant, sliParentId string) ([]*dbtype.Node, error)
	GetServiceLineItemsForContract(ctx context.Context, tenant, contractId string) ([]*dbtype.Node, error)
	GetServiceLineItemsForContracts(ctx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error)
	GetServiceLineItemsForInvoiceLines(ctx context.Context, tenant string, invoiceLineIds []string) ([]*utils.DbNodeAndId, error)
	GetLatestServiceLineItemByParentId(ctx context.Context, tenant, serviceLineItemParentId string, beforeDate *time.Time) (*dbtype.Node, error)
	WasServiceLineItemInvoiced(ctx context.Context, tenant, serviceLineItemId string) (bool, error)
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

func (r *serviceLineItemReadRepository) GetServiceLineItemsForContract(ctx context.Context, tenant, contractId string) ([]*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemReadRepository.GetServiceLineItemsForContract")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("contractId", contractId))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})-[:HAS_SERVICE]->(sli:ServiceLineItem)
							WHERE sli:ServiceLineItem_%s
							RETURN sli ORDER BY sli.createdAt ASC`, tenant)
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

func (r *serviceLineItemReadRepository) GetServiceLineItemsForContracts(ctx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemRepository.GetServiceLineItemsForContracts")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("contractIds", contractIds))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)-[:HAS_SERVICE]->(sli:ServiceLineItem)
			WHERE c.id IN $contractIds and sli:ServiceLineItem_%s
			RETURN sli, c.id ORDER BY sli.createdAt ASC`, tenant)
	params := map[string]any{
		"tenant":      tenant,
		"contractIds": contractIds,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *serviceLineItemReadRepository) GetServiceLineItemsForInvoiceLines(ctx context.Context, tenant string, invoiceLineIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemRepository.GetServiceLineItemsForInvoiceLines")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("invoiceLineIds", invoiceLineIds))

	cypher := fmt.Sprintf(`MATCH (il:InvoiceLine)-[:INVOICED]->(sli:ServiceLineItem)
			WHERE il.id IN $invoiceLineIds and sli:ServiceLineItem_%s
			RETURN sli, il.id ORDER BY sli.createdAt ASC`, tenant)
	params := map[string]any{
		"tenant":         tenant,
		"invoiceLineIds": invoiceLineIds,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *serviceLineItemReadRepository) GetServiceLineItemsByParentId(ctx context.Context, tenant, sliParentId string) ([]*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemReadRepository.GetServiceLineItemsByParentId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("sliParentId", sliParentId))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)-[:HAS_SERVICE]->(sli:ServiceLineItem {parentId:$parentId})
							WHERE sli:ServiceLineItem_%s
							RETURN sli ORDER BY sli.startedAt`, tenant)
	params := map[string]any{
		"tenant":   tenant,
		"parentId": sliParentId,
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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

func (r *serviceLineItemReadRepository) GetLatestServiceLineItemByParentId(ctx context.Context, tenant, serviceLineItemParentId string, beforeDate *time.Time) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemReadRepository.GetLatestServiceLineItemByParentId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("serviceLineItemParentId", serviceLineItemParentId), log.Object("beforeDate", beforeDate))

	params := map[string]any{
		"tenant":   tenant,
		"parentId": serviceLineItemParentId,
	}
	cypher := `MATCH (sli:ServiceLineItem {parentId:$parentId}) `
	if beforeDate != nil {
		cypher += ` WHERE sli.startedAt < $before `
		params["before"] = beforeDate.Add(time.Millisecond * 1)
	}
	cypher += ` RETURN sli ORDER BY sli.startedAt DESC LIMIT 1`

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

func (r *serviceLineItemReadRepository) WasServiceLineItemInvoiced(ctx context.Context, tenant, serviceLineItemId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemReadRepository.WasServiceLineItemInvoiced")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, serviceLineItemId)

	cypher := fmt.Sprintf(`MATCH (sli:ServiceLineItem {id:$id})<-[:INVOICED]-(il:InvoiceLine)--(i:Invoice {dryRun:false}) WHERE sli:ServiceLineItem_%s RETURN count(sli)`, tenant)
	params := map[string]any{
		"id": serviceLineItemId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsType[int64](ctx, queryResult, err)
	})
	if err != nil {
		span.LogFields(log.Bool("result.found", false))
		tracing.TraceErr(span, err)
		return false, err
	}
	if result.(int64) == 0 {
		span.LogFields(log.Bool("result.found", false))
		return false, nil
	}
	span.LogFields(log.Bool("result.found", true))
	return true, nil
}
