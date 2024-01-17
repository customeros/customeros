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

type InvoiceLineWriteRepository interface {
	CreateInvoiceLine(ctx context.Context, tenant, id string, name string, price float64, quantity int64, amount, vat, total float64, createdAt time.Time) error
}

type invoiceLineWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInvoiceLineWriteRepository(driver *neo4j.DriverWithContext, database string) InvoiceLineWriteRepository {
	return &invoiceLineWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *invoiceLineWriteRepository) CreateInvoiceLine(ctx context.Context, tenant, id string, name string, price float64, quantity int64, amount, vat, total float64, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceLineWriteRepository.CreateInvoiceLine")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$id})
							WHERE i:Invoice_%s
							MERGE (i)-[:HAS_INVOICE_LINE]->(il:InvoiceLine {id:randomUUID()})
							ON CREATE SET 
								il:InvoiceLine_%s,
								il.createdAt=$createdAt,
								il.updatedAt=$updatedAt,
								il.name=$name,
								il.price=$price,
								il.quantity=$quantity,
								il.amount=$amount,
								il.vat=$vat,
								il.total=$total
`, tenant, tenant)
	params := map[string]any{
		"tenant":    tenant,
		"id":        id,
		"createdAt": createdAt,
		"updatedAt": createdAt,
		"name":      name,
		"price":     price,
		"quantity":  quantity,
		"amount":    amount,
		"vat":       vat,
		"total":     total,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
