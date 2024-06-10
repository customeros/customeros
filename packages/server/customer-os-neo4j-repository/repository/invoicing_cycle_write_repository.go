package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type InvoicingCycleWriteRepository interface {
	CreateInvoicingCycleType(ctx context.Context, tenant, id, invoicingCycleType, source, appSource string, createdAt time.Time) error
	UpdateInvoicingCycleType(ctx context.Context, tenant, id, invoicingCycleType string) error
}

type invoicingCycleWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInvoicingCycleWriteRepository(driver *neo4j.DriverWithContext, database string) InvoicingCycleWriteRepository {
	return &invoicingCycleWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *invoicingCycleWriteRepository) CreateInvoicingCycleType(ctx context.Context, tenant, id, invoicingCycleType, source, appSource string, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicingCycleWriteRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
							MERGE (t)<-[:INVOICING_CYCLE_BELONGS_TO_TENANT]-(ic:InvoicingCycle {id:$id}) 
							ON CREATE SET 
								ic:InvoicingCycle_%s,
								ic.createdAt=$createdAt,
								ic.updatedAt=datetime(),
								ic.source=$source,
								ic.sourceOfTruth=$sourceOfTruth,
								ic.appSource=$appSource,
								ic.type=$type
							`, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"id":            id,
		"createdAt":     createdAt,
		"source":        source,
		"sourceOfTruth": source,
		"appSource":     appSource,
		"type":          invoicingCycleType,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *invoicingCycleWriteRepository) UpdateInvoicingCycleType(ctx context.Context, tenant, id, invoicingCycleType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicingCycleWriteRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICING_CYCLE_BELONGS_TO_TENANT]-(ic:InvoicingCycle {id:$id}) 
							SET ic.updatedAt=datetime(),
								ic.type=$type`
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
		"type":   invoicingCycleType,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
