package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type ServiceLineItemCreateFields struct {
	IsNewVersionForExistingSLI bool         `json:"isNewVersionForExistingSLI"`
	PreviousQuantity           int64        `json:"previousQuantity"`
	PreviousPrice              float64      `json:"previousPrice"`
	PreviousBilled             string       `json:"previousBilled"`
	SourceFields               model.Source `json:"sourceFields"`
	ContractId                 string       `json:"contractId"`
	ParentId                   string       `json:"parentId"`
	CreatedAt                  time.Time    `json:"createdAt"`
	StartedAt                  time.Time    `json:"startedAt"`
	EndedAt                    *time.Time   `json:"endedAt"`
	Price                      float64      `json:"price"`
	Quantity                   int64        `json:"quantity"`
	Name                       string       `json:"name"`
	Billed                     string       `json:"billed"`
	Comments                   string       `json:"comments"`
	VatRate                    float64      `json:"vatRate"`
	PreviousVatRate            float64      `json:"previousVatRate"`
}

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
	CreateForContract(ctx context.Context, tenant, serviceLineItemId string, data ServiceLineItemCreateFields) error
	Update(ctx context.Context, tenant, serviceLineItemId string, data ServiceLineItemUpdateFields) error
	Delete(ctx context.Context, tenant, serviceLineItemId string) error
	Close(ctx context.Context, tenant, serviceLineItemId string, endedAt time.Time, isCanceled bool) error
	AdjustEndDates(ctx context.Context, tenant, parentId string) error
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

func (r *serviceLineItemWriteRepository) CreateForContract(ctx context.Context, tenant, serviceLineItemId string, data ServiceLineItemCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemWriteRepository.CreateForContract")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
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
								sli.sourceOfTruth=$sourceOfTruth,
								sli.appSource=$appSource,
								sli.name=$name,
								sli.price=$price,
								sli.quantity=$quantity,
								sli.billed=$billed,
								sli.parentId=$parentId,
				                sli.comments=$comments,
								sli.vatRate=toFloat($vatRate)
							`, tenant)
	params := map[string]any{
		"tenant":            tenant,
		"serviceLineItemId": serviceLineItemId,
		"contractId":        data.ContractId,
		"parentId":          data.ParentId,
		"createdAt":         data.CreatedAt,
		"startedAt":         utils.ToDate(data.StartedAt),
		"endedAt":           utils.TimePtrAsAny(utils.ToDatePtr(data.EndedAt)),
		"source":            data.SourceFields.Source,
		"sourceOfTruth":     data.SourceFields.Source,
		"appSource":         data.SourceFields.AppSource,
		"price":             data.Price,
		"quantity":          data.Quantity,
		"name":              data.Name,
		"billed":            data.Billed,
		"comments":          data.Comments,
		"vatRate":           data.VatRate,
	}
	if data.IsNewVersionForExistingSLI {
		cypher += `, sli.previousQuantity=$previousQuantity, sli.previousPrice=$previousPrice, sli.previousBilled=$previousBilled, sli.previousVatRate=toFloat($previousVatRate)`
		params["previousQuantity"] = data.PreviousQuantity
		params["previousPrice"] = data.PreviousPrice
		params["previousBilled"] = data.PreviousBilled
		params["previousVatRate"] = data.PreviousVatRate
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *serviceLineItemWriteRepository) Update(ctx context.Context, tenant, serviceLineItemId string, data ServiceLineItemUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemWriteRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, serviceLineItemId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (sli:ServiceLineItem {id:$serviceLineItemId})
							WHERE sli:ServiceLineItem_%s
							SET 
								sli.name = CASE WHEN sli.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $name ELSE sli.name END,
								sli.price = CASE WHEN sli.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $price ELSE sli.price END,
								sli.quantity = CASE WHEN sli.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $quantity ELSE sli.quantity END,
								sli.billed = CASE WHEN sli.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $billed ELSE sli.billed END,
								sli.vatRate = CASE WHEN sli.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN toFloat($vatRate) ELSE sli.vatRate END,
								sli.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE sli.sourceOfTruth END,
								sli.updatedAt=datetime(),
				                sli.comments=$comments
							`, tenant)
	params := map[string]any{
		"serviceLineItemId": serviceLineItemId,
		"price":             data.Price,
		"quantity":          data.Quantity,
		"vatRate":           data.VatRate,
		"name":              data.Name,
		"billed":            data.Billed,
		"comments":          data.Comments,
		"sourceOfTruth":     data.Source,
		"overwrite":         data.Source == constants.SourceOpenline,
	}
	if data.StartedAt != nil {
		params["startedAt"] = utils.ToDate(*data.StartedAt)
		cypher += `, sli.startedAt = $startedAt`
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *serviceLineItemWriteRepository) Delete(ctx context.Context, tenant, serviceLineItemId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemWriteRepository.Delete")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
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

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *serviceLineItemWriteRepository) Close(ctx context.Context, tenant, serviceLineItemId string, endedAt time.Time, isCanceled bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemWriteRepository.Close")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, serviceLineItemId)
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

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *serviceLineItemWriteRepository) AdjustEndDates(ctx context.Context, tenant, parentId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemWriteRepository.AdjustEndDates")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
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

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
