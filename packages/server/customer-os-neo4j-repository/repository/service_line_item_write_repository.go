package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/customer-os-neo4j-repository/model"
	"github.com/openline-ai/customer-os-neo4j-repository/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
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
	UpdatedAt                  time.Time    `json:"updatedAt"`
	StartedAt                  time.Time    `json:"startedAt"`
	EndedAt                    *time.Time   `json:"endedAt"`
	Price                      float64      `json:"price"`
	Quantity                   int64        `json:"quantity"`
	Name                       string       `json:"name"`
	Billed                     string       `json:"billed"`
	Comments                   string       `json:"comments"`
}

type ServiceLineItemUpdateFields struct {
	Price     float64   `json:"price"`
	Quantity  int64     `json:"quantity"`
	Name      string    `json:"name"`
	Billed    string    `json:"billed"`
	Comments  string    `json:"comments"`
	Source    string    `json:"source"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ServiceLineItemWriteRepository interface {
	CreateForContract(ctx context.Context, tenant, serviceLineItemId string, data ServiceLineItemCreateFields) error
	Update(ctx context.Context, tenant, serviceLineItemId string, data ServiceLineItemUpdateFields) error
	Delete(ctx context.Context, tenant, serviceLineItemId string) error
	Close(ctx context.Context, tenant, serviceLineItemId string, updatedAt, endedAt time.Time, isCanceled bool) error
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
								sli.updatedAt=$updatedAt,
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
				                sli.comments=$comments
							`, tenant)
	params := map[string]any{
		"tenant":            tenant,
		"serviceLineItemId": serviceLineItemId,
		"contractId":        data.ContractId,
		"parentId":          data.ParentId,
		"createdAt":         data.CreatedAt,
		"updatedAt":         data.UpdatedAt,
		"startedAt":         data.StartedAt,
		"endedAt":           utils.TimePtrFirstNonNilNillableAsAny(data.EndedAt),
		"source":            data.SourceFields.Source,
		"sourceOfTruth":     data.SourceFields.Source,
		"appSource":         data.SourceFields.AppSource,
		"price":             data.Price,
		"quantity":          data.Quantity,
		"name":              data.Name,
		"billed":            data.Billed,
		"comments":          data.Comments,
	}
	if data.IsNewVersionForExistingSLI {
		cypher += `, sli.previousQuantity=$previousQuantity, sli.previousPrice=$previousPrice, sli.previousBilled=$previousBilled`
		params["previousQuantity"] = data.PreviousQuantity
		params["previousPrice"] = data.PreviousPrice
		params["previousBilled"] = data.PreviousBilled
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
								sli.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE sli.sourceOfTruth END,
								sli.updatedAt=$updatedAt,
				                sli.comments=$comments
							`, tenant)
	params := map[string]any{
		"serviceLineItemId": serviceLineItemId,
		"updatedAt":         data.UpdatedAt,
		"price":             data.Price,
		"quantity":          data.Quantity,
		"name":              data.Name,
		"billed":            data.Billed,
		"comments":          data.Comments,
		"sourceOfTruth":     data.Source,
		"overwrite":         data.Source == constants.SourceOpenline,
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

	cypher := fmt.Sprintf(`MATCH (sli:ServiceLineItem {id:$serviceLineItemId})
							WHERE sli:ServiceLineItem_%s
							DETACH DELETE sli`, tenant)
	params := map[string]any{
		"serviceLineItemId": serviceLineItemId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *serviceLineItemWriteRepository) Close(ctx context.Context, tenant, serviceLineItemId string, updatedAt, endedAt time.Time, isCanceled bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemWriteRepository.Close")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, serviceLineItemId)
	span.LogFields(log.Object("updatedAt", updatedAt), log.Object("endedAt", endedAt), log.Bool("isCanceled", isCanceled))

	params := map[string]any{
		"serviceLineItemId": serviceLineItemId,
		"updatedAt":         updatedAt,
		"endedAt":           endedAt,
	}
	cypher := fmt.Sprintf(`MATCH (sli:ServiceLineItem {id:$serviceLineItemId})
							WHERE sli:ServiceLineItem_%s SET
							sli.endedAt = $endedAt,
							sli.updatedAt = $updatedAt`, tenant)
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
