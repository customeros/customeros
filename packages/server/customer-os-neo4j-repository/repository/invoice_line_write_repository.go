package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type InvoiceLineCreateFields struct {
	CreatedAt               time.Time       `json:"createdAt"`
	SourceFields            model.Source    `json:"sourceFields"`
	Name                    string          `json:"name"`
	Price                   float64         `json:"price"`
	Quantity                int64           `json:"quantity"`
	Amount                  float64         `json:"amount"`
	VAT                     float64         `json:"vat"`
	TotalAmount             float64         `json:"totalAmount"`
	ServiceLineItemId       string          `json:"serviceLineItemId"`
	ServiceLineItemParentId string          `json:"serviceLineItemParentId"`
	BilledType              enum.BilledType `json:"billedType"`
}

type InvoiceLineWriteRepository interface {
	CreateInvoiceLine(ctx context.Context, tenant, invoiceId, invoiceLineId string, data InvoiceLineCreateFields) error
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

func (r *invoiceLineWriteRepository) CreateInvoiceLine(ctx context.Context, tenant, invoiceId, invoiceLineId string, data InvoiceLineCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceLineWriteRepository.CreateInvoiceLine")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceLineId)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$invoiceId})
							MERGE (i)-[:HAS_INVOICE_LINE]->(il:InvoiceLine {id:$invoiceLineId})
							ON CREATE SET 
								il:InvoiceLine_%s,
								il.createdAt=$createdAt,
								il.updatedAt=datetime(),
								il.name=$name,
								il.price=$price,
								il.quantity=$quantity,
								il.amount=$amount,
								il.vat=$vat,
								il.totalAmount=$totalAmount,
								il.billedType=$billedType,
								il.source=$source,
								il.appSource=$appSource,
								il.sourceOfTruth=$sourceOfTruth,
								il.serviceLineItemId=$serviceLineItemId,
								il.serviceLineItemParentId=$serviceLineItemParentId
							WITH il
							MATCH (sli:ServiceLineItem {id:$serviceLineItemId}) 
							WHERE sli:ServiceLineItem_%s AND $serviceLineItemId <> ""
							MERGE (il)-[:INVOICED]->(sli)`, tenant, tenant)
	params := map[string]any{
		"tenant":                  tenant,
		"invoiceId":               invoiceId,
		"invoiceLineId":           invoiceLineId,
		"createdAt":               data.CreatedAt,
		"name":                    data.Name,
		"price":                   data.Price,
		"quantity":                data.Quantity,
		"amount":                  data.Amount,
		"vat":                     data.VAT,
		"totalAmount":             data.TotalAmount,
		"billedType":              data.BilledType.String(),
		"source":                  data.SourceFields.Source,
		"appSource":               data.SourceFields.AppSource,
		"sourceOfTruth":           data.SourceFields.Source,
		"serviceLineItemId":       data.ServiceLineItemId,
		"serviceLineItemParentId": data.ServiceLineItemParentId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
