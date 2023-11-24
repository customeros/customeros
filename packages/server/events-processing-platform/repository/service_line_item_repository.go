package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ServiceLineItemRepository interface {
	CreateForContract(ctx context.Context, tenant, serviceLineItemId string, evt event.ServiceLineItemCreateEvent) error
	Update(ctx context.Context, tenant, serviceLineItemId string, evt event.ServiceLineItemUpdateEvent) error
	GetAllForContract(ctx context.Context, tenant, contractId string) ([]*neo4j.Node, error)
	Delete(ctx context.Context, tenant, serviceLineItemId string) error
}

type serviceLineItemRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewServiceLineItemRepository(driver *neo4j.DriverWithContext, database string) ServiceLineItemRepository {
	return &serviceLineItemRepository{
		driver:   driver,
		database: database,
	}
}

func (r *serviceLineItemRepository) CreateForContract(ctx context.Context, tenant, serviceLineItemId string, evt event.ServiceLineItemCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemRepository.CreateForContract")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId), log.Object("event", evt))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
							MERGE (c)-[:HAS_SERVICE]->(sli:ServiceLineItem {id:$serviceLineItemId})
							ON CREATE SET 
								sli:ServiceLineItem_%s,
								sli.createdAt=$createdAt,
								sli.updatedAt=$updatedAt,
								sli.source=$source,
								sli.sourceOfTruth=$sourceOfTruth,
								sli.appSource=$appSource,
								sli.name=$name,
								sli.price=$price,
								sli.quantity=$quantity,
								sli.billed=$billed
							`, tenant)
	params := map[string]any{
		"tenant":            tenant,
		"serviceLineItemId": serviceLineItemId,
		"contractId":        evt.ContractId,
		"createdAt":         evt.CreatedAt,
		"updatedAt":         evt.UpdatedAt,
		"source":            helper.GetSource(evt.Source.Source),
		"sourceOfTruth":     helper.GetSourceOfTruth(evt.Source.Source),
		"appSource":         helper.GetAppSource(evt.Source.AppSource),
		"price":             evt.Price,
		"quantity":          evt.Quantity,
		"name":              evt.Name,
		"billed":            evt.Billed,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
}

func (r *serviceLineItemRepository) Update(ctx context.Context, tenant, serviceLineItemId string, evt event.ServiceLineItemUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId), log.Object("event", evt))

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
		"updatedAt":         evt.UpdatedAt,
		"price":             evt.Price,
		"quantity":          evt.Quantity,
		"name":              evt.Name,
		"billed":            evt.Billed,
		"comments":          evt.Comments,
		"sourceOfTruth":     helper.GetSourceOfTruth(evt.Source.Source),
		"overwrite":         helper.GetSourceOfTruth(evt.Source.Source) == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
}

func (r *serviceLineItemRepository) GetAllForContract(ctx context.Context, tenant, contractId string) ([]*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemRepository.GetAllForContract")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("contractId", contractId))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})-[:HAS_SERVICE]->(sli:ServiceLineItem)
							WHERE sli:ServiceLineItem_%s
							RETURN sli`, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*neo4j.Node), nil
}

func (r *serviceLineItemRepository) Delete(ctx context.Context, tenant, serviceLineItemId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemRepository.Delete")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId))

	cypher := fmt.Sprintf(`MATCH (sli:ServiceLineItem {id:$serviceLineItemId})
							WHERE sli:ServiceLineItem_%s
							DETACH DELETE sli`, tenant)
	params := map[string]any{
		"serviceLineItemId": serviceLineItemId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
}
