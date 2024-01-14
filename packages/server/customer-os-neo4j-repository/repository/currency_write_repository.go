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

type CurrencyWriteRepository interface {
	CreateCurrency(ctx context.Context, id, name, symbol, source, appSource string, createdAt time.Time) error
}

type currencyWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCurrencyWriteRepository(driver *neo4j.DriverWithContext, database string) CurrencyWriteRepository {
	return &currencyWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *currencyWriteRepository) CreateCurrency(ctx context.Context, id, name, symbol, source, appSource string, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CurrencyWriteRepository.CreateCurrency")
	defer span.Finish()
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := fmt.Sprintf(`MERGE (c:Currency {id:$id})
							ON CREATE SET 
								c.createdAt=$createdAt,
								c.updatedAt=$updatedAt,
								c.source=$source,
								c.sourceOfTruth=$sourceOfTruth,
								c.appSource=$appSource,
								c.name=$name,
								c.symbol=$symbol
							`)
	params := map[string]any{
		"id":            id,
		"createdAt":     createdAt,
		"updatedAt":     createdAt,
		"source":        source,
		"sourceOfTruth": source,
		"appSource":     appSource,
		"name":          name,
		"symbol":        symbol,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
