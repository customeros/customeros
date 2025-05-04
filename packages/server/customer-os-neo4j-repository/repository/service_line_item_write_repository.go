package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"time"
)

type ServiceLineItemWriteRepository interface {
	CreateForContract(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, serviceLineItemId string, data data_fields.SLIFields) error
	Update(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, serviceLineItemId string, data data_fields.SLIFields) error
	Delete(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, serviceLineItemId string) error
	Close(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, serviceLineItemId string, endedAt time.Time, isCanceled bool) error
	AdjustEndDates(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, parentId string) error
}

type serviceLineItemWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewServiceLineItemWriteRepository(driver *neo4j.DriverWithContext, database string) ServiceLineItemWriteRepository {
	return &serviceLineItemWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *serviceLineItemWriteRepository) CreateForContract(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, serviceLineItemId string, data data_fields.SLIFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ServiceLineItemWriteRepository.CreateForContract")
	defer spans.Finish()

	spans.TagEntity(serviceLineItemId)
	spans.LogObjectAsJson("data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
							MERGE (c)-[:HAS_SERVICE]->(sli:ServiceLineItem {id:$serviceLineItemId})
							ON CREATE SET 
								sli:ServiceLineItem_%s,
								sli.createdAt=$createdAt,
								sli.updatedAt=datetime(),
								sli.startedAt=$startedAt,
								sli.endedAt=$endedAt,
								sli.source=$source,
								sli.appSource=$appSource,
								sli.skuId=$skuId,
								sli.description=$description,
								sli.price=toFloat($price),
								sli.quantity=$quantity,
								sli.billed=$billed,
								sli.parentId=$parentId,
				                sli.comments=$comments,
								sli.vatRate=toFloat($vatRate)
							`, tenant)
	params := map[string]any{
		"tenant":            tenant,
		"serviceLineItemId": serviceLineItemId,
		"createdAt":         utils.IfNotNilTimeWithDefault(data.CreatedAt, utils.Now()),
		"startedAt":         utils.ToDate(utils.IfNotNilTimeWithDefault(data.StartedAt, utils.Now())),
		"endedAt":           utils.TimePtrAsAny(utils.ToDatePtr(data.EndedAt)),
		"source":            utils.IfNotNilString(data.Source),
		"appSource":         utils.IfNotNilString(data.AppSource),
		"contractId":        utils.IfNotNilString(data.ContractId),
		"parentId":          utils.IfNotNilString(data.ParentId),
		"price":             utils.IfNotNilFloat64(data.Price),
		"quantity":          utils.IfNotNilInt64(data.Quantity),
		"skuId":             utils.IfNotNilString(data.SkuId),
		"description":       utils.IfNotNilString(data.Description),
		"comments":          utils.IfNotNilString(data.Comments),
		"vatRate":           utils.IfNotNilFloat64(data.TaxRate),
	}
	if data.BilledType != nil {
		params["billed"] = data.BilledType.String()
	} else {
		params["billed"] = ""
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *serviceLineItemWriteRepository) Update(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, serviceLineItemId string, data data_fields.SLIFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ServiceLineItemWriteRepository.Update")
	defer spans.Finish()

	spans.TagEntity(serviceLineItemId)
	spans.LogObjectAsJson("data", data)

	cypher := fmt.Sprintf(`MATCH (sli:ServiceLineItem_%s {id:$serviceLineItemId})
							SET sli.updatedAt=datetime()`, tenant)
	params := map[string]any{
		"serviceLineItemId": serviceLineItemId,
	}
	if data.SkuId != nil {
		params["skuId"] = *data.SkuId
		cypher += `, sli.skuId = $skuId`
	}
	if data.Description != nil {
		params["description"] = *data.Description
		cypher += `, sli.description = $description`
	}
	if data.Price != nil {
		params["price"] = *data.Price
		cypher += `, sli.price = toFloat($price)`
	}
	if data.Quantity != nil {
		params["quantity"] = *data.Quantity
		cypher += `, sli.quantity = $quantity`
	}
	if data.StartedAt != nil {
		params["startedAt"] = utils.ToDate(*data.StartedAt)
		cypher += `, sli.startedAt = $startedAt`
	}
	if data.BilledType != nil {
		params["billed"] = data.BilledType.String()
		cypher += `, sli.billed = $billed`
	}
	if data.TaxRate != nil {
		params["vatRate"] = *data.TaxRate
		cypher += `, sli.vatRate = toFloat($vatRate)`
	}
	if data.Comments != nil {
		params["comments"] = *data.Comments
		cypher += `, sli.comments = $comments`
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *serviceLineItemWriteRepository) Delete(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, serviceLineItemId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ServiceLineItemWriteRepository.Delete")
	defer spans.Finish()

	spans.TagEntity(serviceLineItemId)

	cypher := `MATCH (sli:ServiceLineItem {id:$serviceLineItemId})<-[:HAS_SERVICE]-(c:Contract)-[:CONTRACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
							WHERE sli:ServiceLineItem
							AND NOT (sli)--(:InvoiceLine)--(:Invoice {dryRun:false})
							SET c.updatedAt = datetime()
							DETACH DELETE sli`
	params := map[string]any{
		"tenant":            tenant,
		"serviceLineItemId": serviceLineItemId,
		"now":               utils.Now(),
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *serviceLineItemWriteRepository) Close(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, serviceLineItemId string, endedAt time.Time, isCanceled bool) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ServiceLineItemWriteRepository.Close")
	defer spans.Finish()

	spans.TagEntity(serviceLineItemId)
	spans.LogObjectAsJson("endedAt", endedAt)
	spans.LogKV("isCanceled", isCanceled)

	params := map[string]any{
		"serviceLineItemId": serviceLineItemId,
		"endedAt":           utils.ToDate(endedAt),
	}
	cypher := fmt.Sprintf(`MATCH (sli:ServiceLineItem {id:$serviceLineItemId})
							WHERE sli:ServiceLineItem_%s SET
							sli.endedAt = $endedAt,
							sli.updatedAt = datetime()`, tenant)
	if isCanceled {
		params["isCanceled"] = isCanceled
		cypher += `, sli.isCanceled = $isCanceled`
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *serviceLineItemWriteRepository) AdjustEndDates(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, parentId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ServiceLineItemWriteRepository.AdjustEndDates")
	defer spans.Finish()

	spans.LogKV("parentId", parentId)

	cypher := fmt.Sprintf(`MATCH (sli:ServiceLineItem {parentId: $parentId})
									WHERE sli:ServiceLineItem_%s
									WITH sli
									ORDER BY sli.startedAt ASC
									WITH collect(sli) AS nodes
									FOREACH(i in RANGE(0, size(nodes)-2) | 
    									FOREACH(node in [nodes[i]] | 
        									FOREACH(nextNode in [nodes[i+1]] | 
            									SET node.endedAt = nextNode.startedAt
        									)
    									)
									)
									WITH nodes[size(nodes)-1] AS lastVersion
									WHERE (lastVersion.isCanceled IS NULL OR lastVersion.isCanceled = false)
									SET lastVersion.endedAt = null`, tenant)
	params := map[string]any{
		"parentId": parentId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	if err != nil {
		spans.TraceError(err)
	}

	return err
}
