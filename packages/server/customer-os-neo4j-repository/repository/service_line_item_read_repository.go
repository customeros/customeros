package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ServiceLineItemReadRepository.GetServiceLineItemsForContract")
	defer spans.Finish()

	spans.LogKV("contractId", contractId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})-[:HAS_SERVICE]->(sli:ServiceLineItem)
							WHERE sli:ServiceLineItem_%s
							RETURN sli ORDER BY sli.createdAt ASC`, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*neo4j.Node)))
	return result.([]*neo4j.Node), nil
}

func (r *serviceLineItemReadRepository) GetServiceLineItemsForContracts(ctx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ServiceLineItemRepository.GetServiceLineItemsForContracts")
	defer spans.Finish()

	spans.LogObjectAsJson("contractIds", contractIds)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)-[:HAS_SERVICE]->(sli:ServiceLineItem)
			WHERE c.id IN $contractIds and sli:ServiceLineItem_%s
			RETURN sli, c.id ORDER BY sli.createdAt ASC`, tenant)
	params := map[string]any{
		"tenant":      tenant,
		"contractIds": contractIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ServiceLineItemRepository.GetServiceLineItemsForInvoiceLines")
	defer spans.Finish()

	spans.LogObjectAsJson("invoiceLineIds", invoiceLineIds)

	cypher := fmt.Sprintf(`MATCH (il:InvoiceLine)-[:INVOICED]->(sli:ServiceLineItem)
			WHERE il.id IN $invoiceLineIds and sli:ServiceLineItem_%s
			RETURN sli, il.id ORDER BY sli.createdAt ASC`, tenant)
	params := map[string]any{
		"tenant":         tenant,
		"invoiceLineIds": invoiceLineIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ServiceLineItemReadRepository.GetServiceLineItemsByParentId")
	defer spans.Finish()

	spans.LogKV("sliParentId", sliParentId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)-[:HAS_SERVICE]->(sli:ServiceLineItem {parentId:$parentId})
							WHERE sli:ServiceLineItem_%s
							RETURN sli ORDER BY sli.startedAt`, tenant)
	params := map[string]any{
		"tenant":   tenant,
		"parentId": sliParentId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*neo4j.Node)))
	return result.([]*neo4j.Node), nil
}

func (r *serviceLineItemReadRepository) GetServiceLineItemById(ctx context.Context, tenant, serviceLineItemId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ServiceLineItemReadRepository.GetServiceLineItemById")
	defer spans.Finish()

	spans.TagEntity(serviceLineItemId)

	cypher := fmt.Sprintf(`MATCH (sli:ServiceLineItem {id:$id}) WHERE sli:ServiceLineItem_%s RETURN sli`, tenant)
	params := map[string]any{
		"id": serviceLineItemId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *serviceLineItemReadRepository) GetLatestServiceLineItemByParentId(ctx context.Context, tenant, serviceLineItemParentId string, beforeDate *time.Time) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ServiceLineItemReadRepository.GetLatestServiceLineItemByParentId")
	defer spans.Finish()

	spans.LogKV("serviceLineItemParentId", serviceLineItemParentId)
	spans.LogObjectAsJson("beforeDate", beforeDate)

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

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *serviceLineItemReadRepository) WasServiceLineItemInvoiced(ctx context.Context, tenant, serviceLineItemId string) (bool, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ServiceLineItemReadRepository.WasServiceLineItemInvoiced")
	defer spans.Finish()

	spans.TagEntity(serviceLineItemId)

	cypher := fmt.Sprintf(`MATCH (sli:ServiceLineItem {id:$id})<-[:INVOICED]-(il:InvoiceLine)--(i:Invoice {dryRun:false}) WHERE sli:ServiceLineItem_%s RETURN count(sli)`, tenant)
	params := map[string]any{
		"id": serviceLineItemId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsType[int64](ctx, queryResult, err)
	})
	if err != nil {
		spans.LogKV("result.found", false)
		spans.TraceError(err)
		return false, err
	}
	if result.(int64) == 0 {
		spans.LogKV("result.found", false)
		return false, nil
	}
	spans.LogKV("result.found", true)
	return true, nil
}
