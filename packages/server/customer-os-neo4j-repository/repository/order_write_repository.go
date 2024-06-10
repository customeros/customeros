package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type OrderWriteRepository interface {
	UpsertOrder(ctx context.Context, tenant, organizationId, orderId string, createdAt time.Time, confirmedAt, paidAt, fulfilledAt, canceledAt *time.Time, sourceFields model.Source) error
}

type orderWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOrderWriteRepository(driver *neo4j.DriverWithContext, database string) OrderWriteRepository {
	return &orderWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *orderWriteRepository) UpsertOrder(ctx context.Context, tenant, organizationId, orderId string, createdAt time.Time, confirmedAt, paidAt, fulfilledAt, canceledAt *time.Time, sourceFields model.Source) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderWriteRepository.UpsertOrder")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, orderId)
	tracing.LogObjectAsJson(span, "createdAt", createdAt)
	tracing.LogObjectAsJson(span, "organizationId", organizationId)
	tracing.LogObjectAsJson(span, "confirmedAt", confirmedAt)
	tracing.LogObjectAsJson(span, "paidAt", paidAt)
	tracing.LogObjectAsJson(span, "fulfilledAt", fulfilledAt)
	tracing.LogObjectAsJson(span, "canceledAt", canceledAt)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$organizationId})
							MERGE (t)<-[:ORDER_BELONGS_TO_TENANT]-(or:Order {id:$orderId}) 
							ON CREATE SET
								or.createdAt=$createdAt,
								or:Order_%s,
								or:TimelineEvent,
								or:TimelineEvent_%s
							SET 
								or.updatedAt=datetime(),
								or.source=$source,
								or.sourceOfTruth=$sourceOfTruth,
								or.appSource=$appSource
							
							`, tenant, tenant)
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"orderId":        orderId,
		"createdAt":      createdAt,
		"source":         sourceFields.Source,
		"sourceOfTruth":  sourceFields.Source,
		"appSource":      sourceFields.AppSource,
	}

	if confirmedAt != nil {
		cypher += `, or.confirmedAt=$confirmedAt`
		params["confirmedAt"] = *confirmedAt
	}
	if paidAt != nil {
		cypher += `, or.paidAt=$paidAt`
		params["paidAt"] = *paidAt
	}
	if fulfilledAt != nil {
		cypher += `, or.fulfilledAt=$fulfilledAt`
		params["fulfilledAt"] = *fulfilledAt
	}
	if canceledAt != nil {
		cypher += `, or.canceledAt=$canceledAt`
		params["canceledAt"] = *canceledAt
	}
	cypher += ` WITH o, or 
				MERGE (o)-[:HAS]->(or) `

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
