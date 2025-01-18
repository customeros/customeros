package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type ServiceLineItemUpdateFields struct {
	Price     float64    `json:"price"`
	Quantity  int64      `json:"quantity"`
	Name      string     `json:"name"`
	Billed    string     `json:"billed"`
	Comments  string     `json:"comments"`
	Source    string     `json:"source"`
	VatRate   float64    `json:"vatRate"`
	StartedAt *time.Time `json:"startedAt"`
}

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemWriteRepository.CreateForContract")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, serviceLineItemId)
	tracing.LogObjectAsJson(span, "data", data)

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
								sli.name=$name,
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
		"name":              utils.IfNotNilString(data.Name),
		"comments":          utils.IfNotNilString(data.Comments),
		"vatRate":           utils.IfNotNilFloat64(data.TaxRate),
	}
	if data.BilledType != nil {
		params["billed"] = data.BilledType.String()
	} else {
		params["billed"] = ""
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *serviceLineItemWriteRepository) Update(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, serviceLineItemId string, data data_fields.SLIFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemWriteRepository.Update")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, serviceLineItemId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (sli:ServiceLineItem_%s {id:$serviceLineItemId})
							SET sli.updatedAt=datetime()`, tenant)
	params := map[string]any{
		"serviceLineItemId": serviceLineItemId,
	}
	if data.Name != nil {
		params["name"] = *data.Name
		cypher += `, sli.name = $name`
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
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *serviceLineItemWriteRepository) Delete(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, serviceLineItemId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemWriteRepository.Delete")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, serviceLineItemId)

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
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *serviceLineItemWriteRepository) Close(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, serviceLineItemId string, endedAt time.Time, isCanceled bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemWriteRepository.Close")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, serviceLineItemId)
	span.LogFields(log.Object("endedAt", endedAt), log.Bool("isCanceled", isCanceled))

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
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *serviceLineItemWriteRepository) AdjustEndDates(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, parentId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemWriteRepository.AdjustEndDates")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag("parentId", parentId)

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
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}
