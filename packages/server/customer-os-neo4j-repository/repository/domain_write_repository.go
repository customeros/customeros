package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type DomainWriteRepository interface {
	MergeDomain(ctx context.Context, domain, source, appSource string, now time.Time) error
	EnrichFailed(ctx context.Context, domain, enrichError string, enrichSource enum.DomainEnrichSource, requestedAt time.Time) error
	EnrichSuccess(ctx context.Context, domain, enrichData string, enrichSource enum.DomainEnrichSource, enrichedAt time.Time) error
}

type domainWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewDomainWriteRepository(driver *neo4j.DriverWithContext, database string) DomainWriteRepository {
	return &domainWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (d domainWriteRepository) MergeDomain(ctx context.Context, domain, source, appSource string, time time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainWriteRepository.MergeDomain")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, domain)
	span.SetTag(tracing.SpanTagEntityId, domain)
	tracing.LogObjectAsJson(span, "data", domain)

	cypher := fmt.Sprintf(`
	MERGE (d:Domain {domain:$domain})
	ON CREATE SET
		d.createdAt=$createdAt,
		d.updatedAt=$updatedAt,
		d.source=$source,
		d.appSource=$appSource
	RETURN d
`)

	params := map[string]interface{}{
		"domain":    domain,
		"createdAt": time,
		"updatedAt": time,
		"source":    source,
		"appSource": appSource,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *d.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (d domainWriteRepository) EnrichFailed(ctx context.Context, domain, enrichError string, enrichSource enum.DomainEnrichSource, requestedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainWriteRepository.EnrichFailed")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, domain)
	span.SetTag(tracing.SpanTagEntityId, domain)
	span.LogFields(log.String("enrichError", enrichError), log.String("requestedAt", requestedAt.String()))

	cypher := fmt.Sprintf(`
	MATCH (d:Domain {domain:$domain})
	SET
		d.enrichError=$enrichError,
		d.enrichSource=$enrichSource,
		d.enrichRequestedAt=$enrichRequestedAt`)

	params := map[string]interface{}{
		"domain":            domain,
		"enrichError":       enrichError,
		"enrichRequestedAt": requestedAt,
		"enrichSource":      enrichSource.String(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *d.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err

}

func (d domainWriteRepository) EnrichSuccess(ctx context.Context, domain, enrichData string, enrichSource enum.DomainEnrichSource, enrichedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainWriteRepository.EnrichSuccess")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, domain)
	span.SetTag(tracing.SpanTagEntityId, domain)
	span.LogFields(log.String("enrichData", enrichData), log.String("enrichedAt", enrichedAt.String()))

	cypher := `MATCH (d:Domain {domain:$domain})
	SET
		d.enrichData=$enrichData,
		d.enrichSource=$enrichSource,
		d.enrichRequestedAt=$enrichRequestedAt,
		d.enrichedAt=$enrichedAt
	REMOVE d.enrichError`

	params := map[string]interface{}{
		"domain":            domain,
		"enrichRequestedAt": enrichedAt,
		"enrichedAt":        enrichedAt,
		"enrichData":        enrichData,
		"enrichSource":      enrichSource.String(),
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *d.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
